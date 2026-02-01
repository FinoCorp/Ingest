# Ingest-CLI
_cli application to turn messy input into a known, validated shape._


## goals first version
- setup cobra cli **DONE**
- basic commands: <br>
    -- 'detect' to print the content of the file.<br>
    -- 'normalize' to output a normalized file (csv format)<br>
    -- 'validate' to output a report (log, problems)<br>
- basic flags:<br>
    -- 'delimeter' set custom delimeter.<br>
    -- 'date-format' set custom date format.<br>
    -- 'out'<br>
    -- 'strict' stop run on first error and collect all log data.<br>
- define standard config/schema<br>
- define normalization rules<br>
- define validation rules<br>
