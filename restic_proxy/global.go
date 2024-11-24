package resticProxy

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hantbk/vbackup/internal/consts/system_status"
	repoModel "github.com/hantbk/vbackup/internal/entity/v1/repository"
	"github.com/hantbk/vbackup/internal/server"
	"github.com/hantbk/vbackup/internal/service/v1/common"
	repositoryDao "github.com/hantbk/vbackup/internal/service/v1/repository"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/azure"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/b2"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/gs"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/limiter"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/local"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/location"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/rclone"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/rest"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/retry"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/s3"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/sftp"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/backend/swift"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/cache"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/debug"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/errors"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/fs"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/options"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/repository"
	"github.com/hantbk/vbackup/pkg/restic_source/rinternal/restic"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

// TimeFormat is the format used for all timestamps printed by restic.
const TimeFormat = "2006-01-02 15:04:05"

var version = "0.16.5"

type backendWrapper func(r restic.Backend) (restic.Backend, error)

// GlobalOptions hold all global options for restic.
type GlobalOptions struct {
	ctx           context.Context
	Repo          string
	KeyHint       string
	Quiet         bool
	Verbose       int
	NoLock        bool
	RetryLock     time.Duration
	JSON          bool
	CacheDir      string
	NoCache       bool
	CleanupCache  bool
	Compression   repository.CompressionMode
	PackSize      uint
	NoExtraVerify bool

	backend.TransportOptions
	limiter.Limits

	// AWS_ACCESS_KEY_ID
	KeyId string
	// AWS_SECRET_ACCESS_KEY
	Secret string
	// AWS_DEFAULT_REGION
	Region string
	// GOOGLE_PROJECT_ID
	ProjectID string
	// AZURE_ACCOUNT_NAME
	AccountName string
	// AZURE_ACCOUNT_KEY，B2_ACCOUNT_KEY
	AccountKey string
	// B2_ACCOUNT_ID
	AccountID string
	password  string

	backends                              *location.Registry
	backendTestHook, backendInnerTestHook backendWrapper

	// verbosity is set as follows:
	//  0 means: don't print any messages except errors, this is used when --quiet is specified
	//  1 is the default: print essential messages
	//  2 means: print more messages, report minor things, this is used when --verbose is specified
	//  3 means: print very detailed debug messages, this is used when --verbose=2 is specified
	verbosity uint

	Options []string

	extended options.Options

	RepositoryVersion string
}

type Repository struct {
	repoId   int
	repoName string
	repo     *repository.Repository
	cancel   context.CancelFunc
	gopts    GlobalOptions
}

var repositoryService repositoryDao.Service

// GetGlobalOptions: Retrieve repository configuration
func GetGlobalOptions(rep repoModel.Repository) (GlobalOptions, context.CancelFunc) {
	var types string
	switch rep.Type {
	case repoModel.S3:
		types = "s3:"
	case repoModel.Sftp:
		types = "sftp:"
	case repoModel.Rest:
		types = "rest:"
	case repoModel.Local:
		types = ""
	default:
		types = ""
	}
	var repo string
	if rep.Type == repoModel.Rest {
		endpoint, err := url.Parse(rep.Endpoint)
		if err != nil {
			server.Logger().Error(err)
			return GlobalOptions{}, nil
		}
		repo = types + endpoint.Scheme + "://" + rep.KeyId + ":" + rep.Secret + "@" + endpoint.Host + endpoint.Path
	} else {
		repo = types + rep.Endpoint + "/" + rep.Bucket
	}
	if rep.PackSize == 0 {
		rep.PackSize = 16
	}
	var globalOptions = GlobalOptions{
		Repo:              repo,
		KeyId:             rep.KeyId,
		Secret:            rep.Secret,
		Region:            rep.Region,
		CleanupCache:      true,
		Compression:       repository.CompressionMode(rep.Compression),
		PackSize:          uint(rep.PackSize),
		NoExtraVerify:     false,
		ProjectID:         rep.ProjectID,
		AccountName:       rep.AccountName,
		AccountKey:        rep.AccountKey,
		AccountID:         rep.AccountID,
		password:          rep.Password,
		RepositoryVersion: rep.RepositoryVersion,
		CacheDir:          server.Config().Data.CacheDir,
		NoCache:           server.Config().Data.NoCache,
		Options:           []string{},
	}
	backends := location.NewRegistry()
	backends.Register(azure.NewFactory())
	backends.Register(b2.NewFactory())
	backends.Register(gs.NewFactory())
	backends.Register(local.NewFactory())
	backends.Register(rclone.NewFactory())
	backends.Register(rest.NewFactory())
	backends.Register(s3.NewFactory())
	backends.Register(sftp.NewFactory())
	backends.Register(swift.NewFactory())

	globalOptions.backends = backends
	var cancel context.CancelFunc
	globalOptions.ctx, cancel = context.WithCancel(context.Background())
	return globalOptions, cancel
}

var repositoryLock sync.Mutex

type RepositoryHandler struct {
	rep  map[int]Repository
	lock sync.Mutex
}

var Myrepositorys RepositoryHandler

func (rh *RepositoryHandler) Get(key int) Repository {
	rh.lock.Lock()
	defer rh.lock.Unlock()
	return Myrepositorys.rep[key]
}

func (rh *RepositoryHandler) Set(key int, rep Repository) {
	rh.lock.Lock()
	defer rh.lock.Unlock()
	Myrepositorys.rep[key] = rep
}

func (rh *RepositoryHandler) Remove(key int) {
	rh.lock.Lock()
	defer rh.lock.Unlock()
	delete(Myrepositorys.rep, key)
}

func cleanCtx() {
	for _, rep := range Myrepositorys.rep {
		rep.cancel()
	}
}

func InitRepository() {
	repositoryLock.Lock()
	defer repositoryLock.Unlock()
	server.UpdateSystemStatus(system_status.Loading)
	defer server.UpdateSystemStatus(system_status.Normal)
	reps, err := repositoryService.List(0, "", common.DBOptions{})
	if err != nil && err.Error() != "not found" {
		fmt.Printf("Repository query failed: %v\n", err)
		return
	}
	cleanCtx()
	ctx := context.Background()
	Myrepositorys = RepositoryHandler{rep: make(map[int]Repository)}
	for _, rep := range reps {
		option, cancel := GetGlobalOptions(rep)
		openRepository, err1 := OpenRepository(ctx, option)
		if err1 != nil {
			fmt.Printf("Repository loading failed: %v\n", err1)
			continue
		}
		repoa := Repository{
			repoId:   rep.Id,
			repoName: rep.Name,
			repo:     openRepository,
			cancel:   cancel,
			gopts:    option,
		}
		_, _ = restic.RemoveAllLocks(ctx, openRepository)
		err = openRepository.LoadIndex(option.ctx, nil)
		if err != nil {
			fmt.Printf("Failed to load index for repository %s: %v\n", rep.Name, err)
			continue
		}
		Myrepositorys.Set(rep.Id, repoa)
	}
	go GetAllRepoStats()
	fmt.Println("Repository loaded successfully! ")

}

// GetRepository: Retrieve repository operation object
func GetRepository(repoid int) (*Repository, error) {
	if repoid <= 0 {
		return nil, errors.Errorf("Repository ID cannot be empty.")
	}
	myrepository := Myrepositorys.rep[repoid]
	if myrepository.repo == nil {
		return nil, fmt.Errorf("Repository does not exist!")
	} else {
		return &myrepository, nil
	}
}

func init() {
	repositoryService = repositoryDao.GetService()
}

func ReadRepo(opts GlobalOptions) (string, error) {
	if opts.Repo == "" {
		return "", errors.Errorf("Please specify repository location (-r or --repository-file)")
	}
	repo := opts.Repo
	return repo, nil
}

var (
	mailLock        sync.Mutex
	checkRepoStatus bool = false
)

func SendEmail(ctx context.Context, to, subject, body string) error {

	// dir, err2 := os.Getwd()
	// if err2 != nil {
	// 	log.Fatal("Error getting current directory:", err2)
	// }
	// fmt.Println("Current directory:", dir)

	// Load environment variables from .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}

	// Get the SMTP credentials and other settings from environment variables
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp-relay.brevo.com"
	smtpPort := 587

	adminEmailAddress := os.Getenv("ADMIN_EMAIL_ADDRESS")
	if smtpUser == "" || smtpPassword == "" || adminEmailAddress == "" {
		log.Fatal("SMTP credentials or email addresses are not set properly.")
		return fmt.Errorf("SMTP credentials or email addresses are not set")
	}

	// Create a new email message
	mailer := gomail.NewMessage()

	// Set the sender and recipient
	mailer.SetHeader("From", adminEmailAddress)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)

	// Set the email body (HTML content)
	mailer.SetBody("text/html", body)

	// Set up the SMTP client with the provided credentials
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)

	// Send the email
	err1 := dialer.DialAndSend(mailer)
	if err1 != nil {
		log.Printf("Failed to send email: %v", err1)
		return err1
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

const maxKeys = 20

// OpenRepository reads the password and opens the repository.
func OpenRepository(ctx context.Context, opts GlobalOptions) (*repository.Repository, error) {
	repo, err := ReadRepo(opts)
	if err != nil {
		return nil, err
	}

	be, err := open(ctx, repo, opts, opts.extended)
	if err != nil {
		return nil, err
	}

	curRetry := 0
	maxRetries := 5

	// Report function inside OpenRepository
	report := func(msg string, err error, d time.Duration) {
		curRetry++
		fmt.Printf("%v returned error, retrying after %v: %v\n", msg, d, err)

		if curRetry >= maxRetries {
			curRetry = 0
			repoPath := "Unknown"

			// Ensure backend is not nil
			if be != nil {
				repoPath = be.Location()
			} else {
				log.Printf("Error: Backend is nil, cannot determine repository path")
			}

			fmt.Printf("Repo at path %s is in error state\n", repoPath)

			// Sử dụng goroutine để gửi email bất đồng bộ
			go func() {
				mailLock.Lock()
				defer mailLock.Unlock()

				// Get the current time
				currentTime := time.Now().Format("2006-01-02 15:04:05")

				// Create HTML content for the email
				emailBody := fmt.Sprintf(`
				<html>
					<body style="font-family: Arial, sans-serif; color: #333; background-color: #f4f4f4; padding: 20px; text-align: center;">
						<!-- Container for the email content -->
						<div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);">
							<!-- Logo -->
							<img src="https://github.com/hantbk/vbackup/blob/main/web/dashboard/src/assets/logo/vbackup-bar.png?raw=true" alt="vBackup Logo" style="max-width: 200px; margin-bottom: 20px;" />

							<!-- Header -->
							<h2 style="color: #0056b3; font-size: 24px; margin-bottom: 20px;">Repository Error Notification</h2>

							<!-- Error Information -->
							<p style="font-size: 16px; margin-bottom: 10px;">
								<strong>Error detected:</strong> A repository issue occurred at path:
								<code style="background-color: #f0f0f0; padding: 5px; border-radius: 4px; font-size: 16px; color: #d9534f;">%s</code>
							</p>
							<p style="font-size: 16px; margin-bottom: 10px;">
								<strong>Time of error:</strong> <span style="font-weight: bold;">%s</span>
							</p>

							<!-- Additional Information -->
							<p style="font-size: 16px; margin-bottom: 20px;">
								Please check the repository status and resolve the issue promptly to prevent further disruptions.
							</p>

							<!-- Footer -->
							<hr style="border: 1px solid #ddd; margin: 20px 0;" />
							<p style="font-style: italic; color: gray; font-size: 14px;">
								This is an automated message. Do not reply to this email.
							</p>
						</div>
					</body>
				</html>`, repoPath, currentTime)

				// Send the email
				if err := SendEmail(ctx, "hantbka@gmail.com",
					"[vBackup] Repository Error Notification",
					emailBody); err != nil {
					log.Printf("Failed to send error email: %v", err)
				}
			}()

		}
	}

	success := func(msg string, retries int) {
		fmt.Printf("%v operation successful after %d retries\n", msg, retries)
		checkRepoStatus = true
	}

	be = retry.New(be, maxRetries, report, success)

	// wrap backend if a test specified a hook
	if opts.backendTestHook != nil {
		be, err = opts.backendTestHook(be)
		if err != nil {
			return nil, err
		}
	}

	s, err := repository.New(be, repository.Options{
		Compression:   opts.Compression,
		PackSize:      opts.PackSize * 1024 * 1024,
		NoExtraVerify: opts.NoExtraVerify,
	})
	if err != nil {
		return nil, errors.Fatal(err.Error())
	}

	err = s.SearchKey(opts.ctx, opts.password, maxKeys, opts.KeyHint)
	if err != nil {
		opts.password = ""
		// Incorrect password, try again
		return nil, errors.Errorf("Incorrect repository password")
	}
	id := s.Config().ID
	if len(id) > 8 {
		id = id[:8]
	}

	fmt.Printf("repository %s opened successfully, password is correct\n", id)

	if opts.NoCache {
		return s, nil
	}

	c, err := cache.New(s.Config().ID, opts.CacheDir)
	if err != nil {
		return s, nil
	}

	if c.Created {
		fmt.Printf("created new cache in %v\n", c.Base)
	}

	// start using the cache
	s.UseCache(c)

	oldCacheDirs, err := cache.Old(c.Base)
	if err != nil {
		fmt.Printf("unable to find old cache directories: %v\n", err)
	}

	// nothing more to do if no old cache dirs could be found
	if len(oldCacheDirs) == 0 {
		return s, nil
	}

	// cleanup old cache dirs if instructed to do so
	if opts.CleanupCache {
		fmt.Printf("removing %d old cache dirs from %v\n", len(oldCacheDirs), c.Base)

		for _, item := range oldCacheDirs {
			dir := filepath.Join(c.Base, item.Name())
			err = fs.RemoveAll(dir)
			if err != nil {
				fmt.Printf("unable to remove %v: %v\n", dir, err)
			}
		}
	} else {
		fmt.Printf("found %d old cache directories in %v, run `restic cache --cleanup` to remove them\n",
			len(oldCacheDirs), c.Base)
	}

	return s, nil
}

// parseConfig: Parse specific parameters for each backend configuration
func parseConfig(loc location.Location, gopts GlobalOptions, opts options.Options) (interface{}, error) {
	// only apply options for a particular backend here
	opts = opts.Extract(loc.Scheme)

	switch loc.Scheme {
	case "local":
		cfg := loc.Config.(*local.Config)
		if err := opts.Apply(loc.Scheme, cfg); err != nil {
			return nil, err
		}
		debug.Log("opening local repository at %#v", cfg)
		return cfg, nil

	case "sftp":
		cfg := loc.Config.(*sftp.Config)
		if err := opts.Apply(loc.Scheme, cfg); err != nil {
			return nil, err
		}

		debug.Log("opening sftp repository at %#v", cfg)
		return cfg, nil

	case "s3":
		cfg := loc.Config.(*s3.Config)
		if cfg.KeyID == "" {
			cfg.KeyID = gopts.KeyId
		}

		if cfg.Secret.String() == "" {
			cfg.Secret = options.NewSecretString(gopts.Secret)
		}

		if cfg.KeyID == "" && cfg.Secret.String() != "" {
			return nil, errors.Fatalf("unable to open S3 backend: Key ID (KeyId) is empty")
		} else if cfg.KeyID != "" && cfg.Secret.String() == "" {
			return nil, errors.Fatalf("unable to open S3 backend: Secret (Secret) is empty")
		}

		if cfg.Region == "" {
			cfg.Region = gopts.Region
		}

		if err := opts.Apply(loc.Scheme, cfg); err != nil {
			return nil, err
		}

		debug.Log("opening s3 repository at %#v", cfg)
		return cfg, nil

	case "gs":
		cfg := loc.Config.(*gs.Config)
		if cfg.ProjectID == "" {
			cfg.ProjectID = gopts.ProjectID
		}

		if err := opts.Apply(loc.Scheme, &cfg); err != nil {
			return nil, err
		}

		debug.Log("opening gs repository at %#v", cfg)
		return cfg, nil

	case "azure":
		cfg := loc.Config.(*azure.Config)
		if cfg.AccountName == "" {
			cfg.AccountName = gopts.AccountName
		}

		if cfg.AccountKey.String() == "" {
			cfg.AccountKey = options.NewSecretString(gopts.AccountKey)
		}

		if err := opts.Apply(loc.Scheme, cfg); err != nil {
			return nil, err
		}

		debug.Log("opening gs repository at %#v", cfg)
		return cfg, nil

	case "swift":
		cfg := loc.Config.(*swift.Config)

		if err := opts.Apply(loc.Scheme, cfg); err != nil {
			return nil, err
		}

		debug.Log("opening swift repository at %#v", cfg)
		return cfg, nil

	case "b2":
		cfg := loc.Config.(*b2.Config)

		if cfg.AccountID == "" {
			cfg.AccountID = gopts.AccountID
		}

		if cfg.AccountID == "" {
			return nil, errors.Fatalf("unable to open B2 backend: Account ID (AccountID) is empty")
		}

		if cfg.Key.String() == "" {
			cfg.Key = options.NewSecretString(gopts.AccountKey)
		}

		if cfg.Key.String() == "" {
			return nil, errors.Fatalf("unable to open B2 backend: Key (AccountKey) is empty")
		}

		if err := opts.Apply(loc.Scheme, cfg); err != nil {
			return nil, err
		}

		debug.Log("opening b2 repository at %#v", cfg)
		return cfg, nil

	case "rest":
		cfg := loc.Config.(*rest.Config)
		if err := opts.Apply(loc.Scheme, cfg); err != nil {
			return nil, err
		}

		debug.Log("opening rest repository at %#v", cfg)
		return cfg, nil
	case "rclone":
		cfg := loc.Config.(*rclone.Config)
		if err := opts.Apply(loc.Scheme, cfg); err != nil {
			return nil, err
		}

		debug.Log("opening rest repository at %#v", cfg)
		return cfg, nil

	}

	return nil, errors.Fatalf("invalid backend: %q", loc.Scheme)
}

// Open the backend specified by a location config.
func open(ctx context.Context, s string, gopts GlobalOptions, opts options.Options) (restic.Backend, error) {
	debug.Log("parsing location %v", location.StripPassword(gopts.backends, s))
	loc, err := location.Parse(gopts.backends, s)
	if err != nil {

		return nil, errors.Fatalf("parsing repository location failed: %v", err)
	}

	var be restic.Backend

	cfg, err := parseConfig(loc, gopts, opts)
	if err != nil {
		return nil, err
	}

	rt, err := backend.Transport(gopts.TransportOptions)
	if err != nil {
		return nil, errors.Fatal(err.Error())
	}

	lim := limiter.NewStaticLimiter(gopts.Limits)
	rt = lim.Transport(rt)

	factory := gopts.backends.Lookup(loc.Scheme)
	if factory == nil {
		return nil, errors.Fatalf("invalid backend: %q", loc.Scheme)
	}

	be, err = factory.Open(ctx, cfg, rt, lim)
	if err != nil {
		return nil, errors.Fatalf("unable to open repository at %v: %v", location.StripPassword(gopts.backends, s), err)
	}

	// wrap backend if a test specified an inner hook
	if gopts.backendInnerTestHook != nil {
		be, err = gopts.backendInnerTestHook(be)
		if err != nil {
			return nil, err
		}
	}

	// check if config is there
	fi, err := be.Stat(gopts.ctx, restic.Handle{Type: restic.ConfigFile})
	if err != nil {
		return nil, errors.Fatalf("unable to open config file: %v\nIs there a repository at the following location?\n%v", err, location.StripPassword(gopts.backends, s))
	}

	if fi.Size == 0 {
		return nil, errors.New("config file has zero size, invalid repository?")
	}

	return be, nil
}

// Create the backend specified by URI.
func create(ctx context.Context, s string, gopts GlobalOptions, opts options.Options) (restic.Backend, error) {
	debug.Log("parsing location %v", s)
	loc, err := location.Parse(gopts.backends, s)
	if err != nil {
		return nil, err
	}

	cfg, err := parseConfig(loc, gopts, opts)
	if err != nil {
		return nil, err
	}

	rt, err := backend.Transport(gopts.TransportOptions)
	if err != nil {
		return nil, errors.Fatal(err.Error())
	}

	factory := gopts.backends.Lookup(loc.Scheme)
	if factory == nil {
		return nil, errors.Fatalf("invalid backend: %q", loc.Scheme)
	}

	be, err := factory.Create(ctx, cfg, rt, nil)
	if err != nil {
		return nil, err
	}

	return be, nil
}
