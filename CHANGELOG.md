## 0.1.0 (2026-05-19)

### Feat

- **config**: add key_path key in yaml config
- **init**: add init command to init or add hermit config to your project
- **age**: add age to encrypt secrets
- **status**: add status command and change services yaml property to notes
- **generate**: add generate command
- **manager**: add secrets manager
- add codebase

### Fix

- **unwrap**: make unwrap command copying also non encrypted files
- **wrap**: wrap command tries now to delete .secrets dir if empty

### Refactor

- refactor a lot of things
- **age**: split age internal lib into multiple files
