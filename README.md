Takes in a list of backup archive names containing the date in the filename
and filters them by given prune rules. Inspired by the prune command
of borg.

The following example will make prunerf to keep 7 daily, 4 weekly and 6
monthly backups. The created date will be parsed using the format given
as arg.

```sh
prunef -d 7 -w 4 -m 6 "example.org_%Y-%m-%d_%H-%M-%S.tar.gz"
```

The output is the list of names to prune.

