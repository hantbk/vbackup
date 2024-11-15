# File change detection
*********************

When vbackup encounters a file that has already been backed up, whether in the
current backup or a previous one, it makes sure the file's content is only
stored once in the repository. To do so, it normally has to scan the entire
content of the file. Because this can be very expensive, vbackup also uses a
change detection rule based on file metadata to determine whether a file is
likely unchanged since a previous backup. If it is, the file is not scanned
again.

The previous backup snapshot, called "parent" snapshot in vbackup terminology,
is determined as follows. By default vbackup groups snapshots by hostname and
backup paths, and then selects the latest snapshot in the group that matches
the current backup. You can change the selection criteria using the
``--group-by`` option, which defaults to ``host,paths``. To select the latest
snapshot with the same paths independent of the hostname, use ``paths``. Or,
to only consider the hostname and tags, use ``host,tags``. Alternatively, it
is possible to manually specify a specific parent snapshot using the
``--parent`` option. Finally, note that one would normally set the
``--group-by`` option for the ``forget`` command to the same value.

Change detection is only performed for regular files (not special files,
symlinks or directories) that have the exact same path as they did in a
previous backup of the same location.  If a file or one of its containing
directories was renamed, it is considered a different file and its entire
contents will be scanned again.

Metadata changes (permissions, ownership, etc.) are always included in the
backup, even if file contents are considered unchanged.

On **Unix** (including Linux and Mac), given that a file lives at the same
location as a file in a previous backup, the following file metadata
attributes have to match for its contents to be presumed unchanged:

* Modification timestamp (mtime).
* Metadata change timestamp (ctime).
* File size.
* Inode number (internal number used to reference a file in a filesystem).

The reason for requiring both mtime and ctime to match is that Unix programs
can freely change mtime (and some do). In such cases, a ctime change may be
the only hint that a file did change.
