# Scripts

This directory contains utility scripts for the Bloombox email validation service.

## Refresh Script

The `refresh` script is used to merge and deduplicate domain lists from multiple text files.

### Usage

```bash
go run refresh/main.go <file1.txt> <file2.txt> <output.txt>
```

### Example

```bash
go run refresh/main.go data/disposable.txt data/free.txt merged_domains.txt
```

### What it does

1. **Reads** domain lists from two input files
2. **Deduplicates** domains (removes duplicates across both files)
3. **Sorts** the domains alphabetically
4. **Writes** the unique, sorted domains to the output file

### Input Format

- Input files should contain one domain per line
- Empty lines are ignored
- Leading and trailing whitespace is trimmed

### Output Format

- One domain per line
- Sorted alphabetically
- No duplicates
- No empty lines

### Error Handling

The script will exit with an error code if:

- Incorrect number of arguments is provided
- Input files cannot be read
- Output file cannot be created or written to

### Example Use Cases

- Merging disposable email domain lists
- Combining free email provider lists
- Consolidating blacklist domains
- Creating unified domain lists for validation
