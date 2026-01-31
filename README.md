# Ingest-CLI
cli application to turn messy input into a known, validated shape.

## goals first version
- setup cobra cli **DONE**
- basic commands:
    -- 'detect' to print the content of the file.
    -- 'normalize' to output a normalized file (csv format)
    -- 'validate' to output a report (log, problems)
- basic flags:
    -- 'delimeter' set custom delimeter.
    -- 'date-format' set custom date format.
    -- 'out'
    -- 'strict' stop run on first error and collect all log data.
- define standard config/schema
- define normalization rules
- define validation rules
