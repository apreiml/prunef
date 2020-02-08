Takes in a list of backup archive names containing the date in the filename
and filters them by given prune rules. Inspired by the prune command
of borg.

The following example will make prunerf to keep 7 daily, 4 weekly and 6
monthly backups.

```sh
prunerf \
    --keep-daily    7               \
    --keep-weekly   4               \
    --keep-monthly  6               \
```

The `invert` option will display those backup archives that can be 
removed

```sh
prunerf --keep-daily 7 --invert
```
