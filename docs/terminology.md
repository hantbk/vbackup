Terminology
===========

- *Repository*: All data produced during a backup is sent to and stored in
a repository in a structured form, for example in a file system
hierarchy with several subdirectories. A repository implementation must
be able to fulfill a number of operations, e.g. list the contents.

- *Blob*: A Blob combines a number of data bytes with identifying
information like the SHA-256 hash of the data and its length.

- *Pack*: A Pack combines one or more Blobs, e.g. in a single file.

- *Snapshot*: A Snapshot stands for the state of a file or directory that
has been backed up at some point in time. The state here means the
content and meta data like the name and modification time for the file
or the directory and its contents.

- *Storage ID*: A storage ID is the SHA-256 hash of the content stored in
the repository. This ID is required in order to load the file from the
repository.

Repository Format
=================

All data is stored in a vbackup repository. A repository is able to store
data of several different types, which can later be requested based on
an ID. This so-called "storage ID" is the SHA-256 hash of the content of
a file. All files in a repository are only written once and never
modified afterwards. Writing should occur atomically to prevent concurrent
operations from reading incomplete files. This allows accessing and even
writing to the repository with multiple clients in parallel. Only the ``prune``
operation removes data from the repository.

Repositories consist of several directories and a top-level file called
``config``. For all other files stored in the repository, the name for
the file is the lower case hexadecimal representation of the storage ID,
which is the SHA-256 hash of the file's contents. This allows for easy
verification of files for accidental modifications, like disk read
errors, by simply running the program ``sha256sum`` on the file and
comparing its output to the file name. If the prefix of a filename is
unique amongst all the other files in the same directory, the prefix may
be used instead of the complete filename.

Apart from the files stored within the ``keys`` and ``data`` directories,
all files are encrypted with AES-256 in counter mode (CTR). The integrity
of the encrypted data is secured by a Poly1305-AES message authentication
code (MAC).
Files in the ``data`` directory ("pack files") consist of multiple parts
which are all independently encrypted and authenticated, see below.

In the first 16 bytes of each encrypted file the initialisation vector
(IV) is stored. It is followed by the encrypted data and completed by
the 16 byte MAC. The format is: ``IV || CIPHERTEXT || MAC``. The
complete encryption overhead is 32 bytes. For each file, a new random IV
is selected.

The file ``config`` is encrypted this way and contains a JSON document
like the following:

 ```json

    {
      "version": 2,
      "id": "5956a3f67a6230d4a92cefb29529f10196c7d92582ec305fd71ff6d331d6271b",
      "chunker_polynomial": "25b468838dcb75"
    }
```

After decryption, restic first checks that the version field contains a
version number that it understands, otherwise it aborts. At the moment, the
version is expected to be 1 or 2. The list of changes in the repository
format is contained in the section "Changes" below.

The field ``id`` holds a unique ID which consists of 32 random bytes, encoded
in hexadecimal. This uniquely identifies the repository, regardless if it is
accessed via a remote storage backend or locally. The field
``chunker_polynomial`` contains a parameter that is used for splitting large
files into smaller chunks (see below).

Repository Layout
-----------------

The ``local`` and ``sftp`` backends are implemented using files and
directories stored in a file system. The directory layout is the same
for both backend types and is also used for all other remote backends.

The basic layout of a repository is shown here:

```plaintext

    /tmp/restic-repo
    ├── config
    ├── data
    │   ├── 21
    │   │   └── 2159dd48f8a24f33c307b750592773f8b71ff8d11452132a7b2e2a6a01611be1
    │   ├── 32
    │   │   └── 32ea976bc30771cebad8285cd99120ac8786f9ffd42141d452458089985043a5
    │   ├── 59
    │   │   └── 59fe4bcde59bd6222eba87795e35a90d82cd2f138a27b6835032b7b58173a426
    │   ├── 73
    │   │   └── 73d04e6125cf3c28a299cc2f3cca3b78ceac396e4fcf9575e34536b26782413c
    │   [...]
    ├── index
    │   ├── c38f5fb68307c6a3e3aa945d556e325dc38f5fb68307c6a3e3aa945d556e325d
    │   └── ca171b1b7394d90d330b265d90f506f9984043b342525f019788f97e745c71fd
    ├── keys
    │   └── b02de829beeb3c01a63e6b25cbd421a98fef144f03b9a02e46eff9e2ca3f0bd7
    ├── locks
    ├── snapshots
    │   └── 22a5af1bdc6e616f8a29579458c49627e01b32210d09adb288d1ecda7c5711ec
    └── tmp
```

S3 Legacy Layout (deprecated)
-----------------------------

Unfortunately during development the Amazon S3 backend uses slightly different
paths (directory names use singular instead of plural for ``key``,
``lock``, and ``snapshot`` files), and the pack files are stored directly below
the ``data`` directory. The S3 Legacy repository layout looks like this:


```plaintext
    /config
    /data
     ├── 2159dd48f8a24f33c307b750592773f8b71ff8d11452132a7b2e2a6a01611be1
     ├── 32ea976bc30771cebad8285cd99120ac8786f9ffd42141d452458089985043a5
     ├── 59fe4bcde59bd6222eba87795e35a90d82cd2f138a27b6835032b7b58173a426
     ├── 73d04e6125cf3c28a299cc2f3cca3b78ceac396e4fcf9575e34536b26782413c
    [...]
    /index
     ├── c38f5fb68307c6a3e3aa945d556e325dc38f5fb68307c6a3e3aa945d556e325d
     └── ca171b1b7394d90d330b265d90f506f9984043b342525f019788f97e745c71fd
    /key
     └── b02de829beeb3c01a63e6b25cbd421a98fef144f03b9a02e46eff9e2ca3f0bd7
    /lock
    /snapshot
     └── 22a5af1bdc6e616f8a29579458c49627e01b32210d09adb288d1ecda7c5711ec
```

Pack Format
===========

All files in the repository except Key and Pack files just contain raw
data, stored as ``IV || Ciphertext || MAC``. Pack files may contain one
or more Blobs of data.

A Pack's structure is as follows:


    EncryptedBlob1 || ... || EncryptedBlobN || EncryptedHeader || Header_Length

At the end of the Pack file is a header, which describes the content.
The header is encrypted and authenticated. ``Header_Length`` is the
length of the encrypted header encoded as a four byte integer in
little-endian encoding. Placing the header at the end of a file allows
writing the blobs in a continuous stream as soon as they are read during
the backup phase. This reduces code complexity and avoids having to
re-write a file once the pack is complete and the content and length of
the header is known.

All the blobs (``EncryptedBlob1``, ``EncryptedBlobN`` etc.) are
authenticated and encrypted independently. This enables repository
reorganisation without having to touch the encrypted Blobs. In addition
it also allows efficient indexing, for only the header needs to be read
in order to find out which Blobs are contained in the Pack. Since the
header is authenticated, authenticity of the header can be checked
without having to read the complete Pack.

After decryption, a Pack's header consists of the following elements:


    Type_Blob1 || Data_Blob1 ||
    [...]
    Type_BlobN || Data_BlobN ||

The Blob type field is a single byte. What follows it depends on the type. The following Blob types are defined:

| Type | Meaning              | Data                                                                       |
|------|-----------------------|----------------------------------------------------------------------------|
| 0b00 | data blob            | `Length(encrypted_blob) || Hash(plaintext_blob)`                           |
| 0b01 | tree blob            | `Length(encrypted_blob) || Hash(plaintext_blob)`                           |
| 0b10 | compressed data blob | `Length(encrypted_blob) || Length(plaintext_blob) || Hash(plaintext_blob)` |
| 0b11 | compressed tree blob | `Length(encrypted_blob) || Length(plaintext_blob) || Hash(plaintext_blob)` |


This is enough to calculate the offsets for all the Blobs in the Pack.
The length fields are encoded as four byte integers in little-endian
format. In the Data column, ``Length(plaintext_blob)`` means the length
of the decrypted and uncompressed data a blob consists of.

All other types are invalid, more types may be added in the future. The
compressed types are only valid for repository format version 2. Data and
tree blobs may be compressed with the zstandard compression algorithm.

In repository format version 1, data and tree blobs should be stored in
separate pack files. In version 2, they must be stored in separate files.
Compressed and non-compress blobs of the same type may be mixed in a pack
file.

For reconstructing the index or parsing a pack without an index, first
the last four bytes must be read in order to find the length of the
header. Afterwards, the header can be read and parsed, which yields all
plaintext hashes, types, offsets and lengths of all included blobs.

Unpacked Data Format
====================

Individual files for the index, locks or snapshots are encrypted
and authenticated like Data and Tree Blobs, so the outer structure is
``IV || Ciphertext || MAC`` again. In repository format version 1 the
plaintext always consists of a JSON document which must either be an
object or an array. The JSON encoder must deterministically encode the
document and should match the behavior of the Go standard library implementation
in ``encoding/json``.

Repository format version 2 adds support for compression. The plaintext
now starts with a header to indicate the encoding version to distinguish
it from plain JSON and to allow for further evolution of the storage format:
``encoding_version || data``
The ``encoding_version`` field is encoded as one byte.
For backwards compatibility the encoding versions '[' (0x5b) and '{' (0x7b)
are used to mark that the whole plaintext (including the encoding version
byte) should treated as JSON document.

For new data the encoding version is currently always ``2``. For that
version ``data`` contains a JSON document compressed using the zstandard
compression algorithm.

Indexing
========

Index files contain information about Data and Tree Blobs and the Packs
they are contained in and store this information in the repository. When
the local cached index is not accessible any more, the index files can
be downloaded and used to reconstruct the index. The file encoding is
described in the "Unpacked Data Format" section. The plaintext consists
of a JSON document like the following:

``` javascript

    {
      "supersedes": [
        "ed54ae36197f4745ebc4b54d10e0f623eaaaedd03013eb7ae90df881b7781452"
      ],
      "packs": [
        {
          "id": "73d04e6125cf3c28a299cc2f3cca3b78ceac396e4fcf9575e34536b26782413c",
          "blobs": [
            {
              "id": "3ec79977ef0cf5de7b08cd12b874cd0f62bbaf7f07f3497a5b1bbcc8cb39b1ce",
              "type": "data",
              "offset": 0,
              "length": 38,
              // no 'uncompressed_length' as blob is not compressed
            },
            {
              "id": "9ccb846e60d90d4eb915848add7aa7ea1e4bbabfc60e573db9f7bfb2789afbae",
              "type": "data",
              "offset": 38,
              "length": 112,
              "uncompressed_length": 511,
            },
            {
              "id": "d3dc577b4ffd38cc4b32122cabf8655a0223ed22edfd93b353dc0c3f2b0fdf66",
              "type": "data",
              "offset": 150,
              "length": 123,
              "uncompressed_length": 234,
            }
          ]
        }, [...]
      ]
    }
```

This JSON document lists Packs and the blobs contained therein. In this
example, the Pack ``73d04e61`` contains three data Blobs,
the plaintext hashes are listed afterwards. The ``length`` field
corresponds to ``Length(encrypted_blob)`` in the pack file header.
Field ``uncompressed_length`` is only present for compressed blobs and
therefore is never present in version 1 of the repository format. It is
set to the value of ``Length(blob)``.

The field ``supersedes`` lists the storage IDs of index files that have
been replaced with the current index file. This happens when index files
are repacked, for example when old snapshots are removed and Packs are
recombined.

There may be an arbitrary number of index files, containing information
on non-disjoint sets of Packs. The number of packs described in a single
file is chosen so that the file size is kept below 8 MiB.

Keys, Encryption and MAC
========================

All data stored by restic in the repository is encrypted with AES-256 in
counter mode and authenticated using Poly1305-AES. For encrypting new
data first 16 bytes are read from a cryptographically secure
pseudo-random number generator as a random nonce. This is used both as
the IV for counter mode and the nonce for Poly1305. This operation needs
three keys: A 32 byte for AES-256 for encryption, a 16 byte AES key and
a 16 byte key for Poly1305. For details see the original paper `The
Poly1305-AES message-authentication
code [Link](https://cr.yp.to/mac/poly1305-20050329.pdf)__ by Dan Bernstein.
The data is then encrypted with AES-256 and afterwards a message
authentication code (MAC) is computed over the ciphertext, everything is
then stored as IV \|\| CIPHERTEXT \|\| MAC.

The directory ``keys`` contains key files. These are simple JSON
documents which contain all data that is needed to derive the
repository's master encryption and message authentication keys from a
user's password. The JSON document from the repository can be
pretty-printed for example by using the Python module ``json``
(shortened to increase readability):

```

    $ python -mjson.tool /tmp/restic-repo/keys/b02de82*
    {
        "hostname": "kasimir",
        "username": "fd0"
        "kdf": "scrypt",
        "N": 65536,
        "r": 8,
        "p": 1,
        "created": "2015-01-02T18:10:13.48307196+01:00",
        "data": "tGwYeKoM0C4j4/9DFrVEmMGAldvEn/+iKC3te/QE/6ox/V4qz58FUOgMa0Bb1cIJ6asrypCx/Ti/pRXCPHLDkIJbNYd2ybC+fLhFIJVLCvkMS+trdywsUkglUbTbi+7+Ldsul5jpAj9vTZ25ajDc+4FKtWEcCWL5ICAOoTAxnPgT+Lh8ByGQBH6KbdWabqamLzTRWxePFoYuxa7yXgmj9A==",
        "salt": "uW4fEI1+IOzj7ED9mVor+yTSJFd68DGlGOeLgJELYsTU5ikhG/83/+jGd4KKAaQdSrsfzrdOhAMftTSih5Ux6w==",
    }
```

When the repository is opened by restic, the user is prompted for the
repository password. This is then used with ``scrypt``, a key derivation
function (KDF), and the supplied parameters (``N``, ``r``, ``p`` and
``salt``) to derive 64 key bytes. The first 32 bytes are used as the
encryption key (for AES-256) and the last 32 bytes are used as the
message authentication key (for Poly1305-AES). These last 32 bytes are
divided into a 16 byte AES key ``k`` followed by 16 bytes of secret key
``r``. The key ``r`` is then masked for use with Poly1305 (see the paper
for details).

Those keys are used to authenticate and decrypt the bytes contained in
the JSON field ``data`` with AES-256 and Poly1305-AES as if they were
any other blob (after removing the Base64 encoding). If the
password is incorrect or the key file has been tampered with, the
computed MAC will not match the last 16 bytes of the data, and restic
exits with an error. Otherwise, the data yields a JSON document
which contains the master encryption and message authentication keys for
this repository (encoded in Base64). The command
``restic cat masterkey`` can be used as follows to decrypt and
pretty-print the master key:


All data in the repository is encrypted and authenticated with these
master keys. For encryption, the AES-256 algorithm in Counter mode is
used. For message authentication, Poly1305-AES is used as described
above.

A repository can have several different passwords, with a key file for
each. This way, the password can be changed without having to re-encrypt
all data.
