# prunef

Takes in a list of backup archive names containing the date in the filename
and filters entries that need to be pruned by given rules. The date will be
parsed using a format string similar to the posix format stirngs used in the
date(1) util.

The following example will make prunef to keep 7 daily, 4 weekly and 6
monthly backups. The created date will be parsed using the format given
as arg.

```sh
prunef -keep-daily 7 \
       -keep-weekly 4 \
       -keep-monthly 6 \
       "example.org_%Y-%m-%d_%H-%M-%S.tar.gz"
```

To list backups that will should be kept use the `-inverse` option.

By default it will parse the dates as created from the local timezone
and it works internal with UTC only. If the entries are created in
UTC use the `-utc` flag.

## Backup slots

// TODO

