# pdf-unlocker

Removes password protection from PDFs using [pdfcpu](https://github.com/pdfcpu/pdfcpu).

## Installation

```bash
git clone <repo>
cd pdf-unlocker
go build -o pdf-unlocker.exe .
```

## Usage

```bash
# Single file
go run unlocker.go -in=file.pdf -password=secret [-out=output.pdf]

# Entire directory (processes all PDFs)
go run unlocker.go -in=folder/ -password=secret
```

## Flags

| Flag | Required | Description |
|---|---|---|
| `-in` | yes | Path to a PDF or directory |
| `-password` | yes | PDF password |
| `-out` | no | Output path (single-file only; default: `<name>_unlocked.pdf`) |
