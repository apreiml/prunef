# prunef

Takes an unsorted list of backup names, containing the time created, and
returns a list of backups for deletion. The backup rotation rules are given via
command line args. A date(1) like format specifier is required to parse the
date from the backup names.

The following example will make prunef to keep 7 daily, 4 weekly and 6 monthly
backups.

```sh
prunef --keep-daily 7 \
       --keep-weekly 4 \
       --keep-monthly 6 \
       "example.org_%Y-%m-%d_%H-%M-%S.tar.gz"
```

By default it will parse the dates as created from the local timezone and it
works internal with UTC only. If the entries are created in
UTC use the `--utc` flag.

To list backups that will should be kept use the `--invert` option.

## Backup slots

The algorithm works by creating timestamp slots according to the rules given.
The first slot is always now. A slot ends where the next timestamp begins.
Those are filled with entries, that are between slot timestamp and next
timestamp. If an entry exists, it will be replaced if another matches the
slot and is more recent.

Slots may not be filled, if there is no entry matching the timestamps.
That means you will only fill all slots, if the backups are created
with lower intervals than the lowest slot interval. If for example
a backup is created every 2 hours and `--keep-hourly 4` is provided,
two backups will be selected for pruning.

If no arguments are given, only one slot with now as timestamp will be
created. Hence the latest backup will never be returned to be pruned.

