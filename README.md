Initial commit
<hr>

`prompt` - holds the question message

`type` - can holds one of following values: 
- `confirm` - simple Y/N answer
- `input` - user input 
- `password` - hidden user input (all symbols are masked with `*`) 
- `select` - select one value from list
- `file` - path to file (validate: provided path is file)
- `dir` - path to dir (validate: provided path is directory)
- `private_key` - path to private key (validate: the file is a private key) 
- `public_key` - path to public key (validate: the file is a public key)

`default` - containts the default value for `input` field

`target` - contains the name of the key to be assigned the value

`validate` - contains the field validator which is prevents wrong user input (CURRENTLY IN DEVELOP)

`help` - help message for prompt

Example surveyfile

```yaml
questions:
  - prompt: What is your first name?
    type: input
    default: Vasily
    target: user_first_name

  - prompt: What is your last name?
    type: input
    validate: required
    target: user_last_name

  - prompt: What is your password?
    type: password
    validate: password
    target: user_password

  - prompt: Select ONE option from list
    type: select
    options:
      - One
      - Two
      - Three
    target: user_select

  - prompt: Select file
    type: file
    target: user_file

  - prompt: Select directory
    type: dir
    target: user_dir

  - prompt: Select public key
    type: public_key
    target: ssh_pub_key
    
  - prompt: Select private key
    type: private_key
    target: ssh_priv_key

  - prompt: Like the survey?
    type: confirm
    default: true
    target: user_like
```

example call
```bash
aa-survey -survey Surveyfile.yaml -output result.yaml
```

TODO:
- Complete README.md
- Add Cobra Command https://github.com/spf13/cobra