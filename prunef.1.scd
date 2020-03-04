prunef(1)

# NAME

prunef - a backup rotation filter

# SYNOPSIS

_prunef_ [OPTION]... FORMAT


# DESCRIPTION

Takes an unsorted list of backup names and returns a list of backups for
deletion. The backup rotation rules are given via command line args. The
backup names need to contain the time and a date(1) like format specifier
is required to parse those.

Following OPTIONS are possible:

*--list-kept*
	Show entries to keep instead of entries to prune.

*--keep-daily uint*
	Number of daily entries to keep.

*--keep-hourly uint*
	Number of hourly entries to keep.

*--keep-minutely uint*
	Number of minutely entries to keep.

*--keep-monthly uint*
	Number of monthly entries to keep.

*--keep-secondly uint*
	Number of secondly entries to keep.

*--keep-weekly uint*
	Number of weekly entries to keep.

*--keep-yearly uint*
	Number of yearly entries to keep.

*--print-slots*
	Print slots and exit.

*--utc*
	Expect input dates in UTC.

# EXAMPLE

The following example will make prunef keep 7 daily, 4 weekly and 6 monthly backups.

```
prunef --keep-daily 7 \
       --keep-weekly 4 \
       --keep-monthly 6 \
       "example.org_%Y-%m-%d_%H-%M-%S.tar.gz"
```

# ALGORITHM

The algorithm works by creating timestamp slots according to the rules
given. The first slot is always now. A slot ends where the next timestamp
begins. The slots are filled with entries, that are between a slot timestamp
and the following timestamp. If an entry exists, it will be replaced if
another matches the slot and is more recent.

Slots may not be filled, if there is no entry matching the timestamps. That
means you will only fill all slots, if the backups are created with lower
intervals than the lowest slot interval. If for example a backup is created
every 2 hours and --keep-hourly 4 is provided, only 2 backups will be kept.

If no arguments are given, only one slot with now as timestamp will be
created. Hence the latest backup will never be returned to be pruned.

# AUTHOR

Written by Armin Preiml
