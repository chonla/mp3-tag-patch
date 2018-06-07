# Patch mp3 file with custom ID3 tag

## Create ID3 tag file list

```
mp3-tag-patch list <path>
```

## Patch mp3 file with custom list

```
mp3-tag-patch patch <path>
```

## Use-case

* `mp3-tag-patch list mp3/` to generate `mp3-tags.json` in `mp3` directory.
* Edit `mp3/mp3-tags.json` by put the mp3 information into the corresponding fields.
* Patch the information back to the mp3 files by `mp3-tag-patch patch mp3/`.