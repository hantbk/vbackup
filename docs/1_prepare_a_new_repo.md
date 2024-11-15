
# Preparing a new repository

The place where your backups will be saved is called a "repository". This is simply a directory containing a set of subdirectories and files created by vbackup to store your backups, some corresponding metadata and encryption keys.

To access the repository, a password (also called a key) must be specified. A repository can hold multiple keys that can all be used to access the repository.

Explains how to create ("init") such a repository. The repository can be stored locally, or on some remote server or service. We'll first cover using a local repository; the remaining sections of this chapter cover all the other options. You can skip to the next chapter once you've read the relevant section here.

For automated backups, vbackup supports specifying the repository location in the environment variable ``VBACKUP_REPOSITORY``. vBackup can also read the repository location from a file specified via the ``--repository-file`` option or the environment variable ``VBACKUP_REPOSITORY_FILE``.

For automating the supply of the repository password to vbackup, several options exist:

* Setting the environment variable ``VBACKUP_PASSWORD``

* Specifying the path to a file with the password via the option ``--password-file`` or the environment variable ``VBACKUP_PASSWORD_FILE``

* Configuring a program to be called when the password is needed via the option ``--password-command`` or the environment variable ``VBACKUP_PASSWORD_COMMAND``

## Local

In order to create a repository at ``/srv/vbackup-repo``


Remembering your password is important! If you lose it, you won't be able to access data stored in the repository.

Common Internet File System (CIFS), an implementation of the Server Message Block (SMB) protocol, is used to share file systems, printers, or serial ports over a network

On Linux, storing the backup repository on a CIFS (SMB) share or backing up data from a CIFS share is not recommended due to compatibility issues in older Linux kernels. Either use another backend or set the environment variable `GODEBUG` to `asyncpreemptoff=1`.

## Amazon S3

vBackup can backup data to any Amazon S3 bucket. However, in this case,
changing the URL scheme is not enough since Amazon uses special security
credentials to sign HTTP requests. By consequence, you must first setup
the following environment variables with the credentials you obtained
while creating the bucket.

```bash
export AWS_ACCESS_KEY_ID=<MY_ACCESS_KEY>
export AWS_SECRET_ACCESS_KEY=<MY_SECRET_ACCESS_KEY>
```

When using temporary credentials make sure to include the session token via
the environment variable ``AWS_SESSION_TOKEN``.

You can then easily initialize a repository that uses your Amazon S3 as
a backend. Make sure to use the endpoint for the correct region. The example
uses ``us-east-1``. If the bucket does not exist it will be created in that region:

``s3:s3.us-east-1.amazonaws.com/bucket_name``

vbackup expects path-style URLs for Amazon S3 storage. For example, use ``s3.us-west-2.amazonaws.com/bucket_name``. Virtual-hosted–style URLs like ``bucket_name.s3.us-west-2.amazonaws.com`` are not supported and must be converted to path-style URLs. For more details, refer to the [AWS S3 documentation](https://docs.aws.amazon.com/AmazonS3/latest/userguide/access-bucket-intro.html).


## Minio Server
Minio is an open-source object storage compatible with Amazon S3 API, written in Go.

Download Minio.
Refer to the Minio documentation for installation and setup of Minio Client and Server.
Set up the following environment variables with your Minio Server credentials:

```bash
export AWS_ACCESS_KEY_ID=<YOUR-MINIO-ACCESS-KEY-ID>
export AWS_SECRET_ACCESS_KEY=<YOUR-MINIO-SECRET-ACCESS-KEY>
```
Initialize restic to use Minio as a backend:
```bash
 s3:http://localhost:9000/vbackup init
```
When prompted, enter a password for the new repository. This password is required for future access, so keep it secure.

### S3-Compatible Storage
For S3-compatible services that are not Amazon, specify the server URL as shown:

```bash
s3:https://server:port/bucket_name
```
Set your credentials for authentication:

```bash
export AWS_ACCESS_KEY_ID=<YOUR-ACCESS-KEY-ID>
export AWS_SECRET_ACCESS_KEY=<YOUR-SECRET-ACCESS-KEY>
-r s3:https://server:port/bucket_name init
```
If necessary, set the region using either the ``AWS_DEFAULT_REGION`` environment variable or by passing the -o option to vbackup like ``-o s3.region="us-east-1"``
By default, ``us-east-1`` is used if no region is specified.

### Path-Style vs. Virtual-Hosted Access

To choose between path-style and virtual-hosted access, use the extended option ``-o s3.bucket-lookup=auto``:

- auto: Default. Uses dns for Amazon and Google endpoints and path for others.
- dns: Use virtual-hosted–style bucket access.
- path: Use path-style bucket access.
  
For certain S3-compatible services, like Ceph versions before v14.2.5, that lack support for ``ListObjectsV2``, enable compatibility by adding ``-o s3.list-objects-v1=true``. Note that this workaround may be removed in future releases.

# Future Work

## Backblaze B2
## Google Cloud Storage
## Microsoft Azure Blob Storage
## Wasabi Hot Cloud Storage
## OpenStack Swift
## SFTP
## REST Server
## Other Services via Rclone