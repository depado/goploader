## Features

- Configure the domain of the goploader instance you want to use
- Take a screenshot and immediately send it
- See a progress bar for long uploads (with speed, ETA, etc)
- Tee what's being uploaded, meaning it will also display it to stdout
- Define the lifetime of your file easily
- Define if your file can be downloaded only once
- Copy the resulting URL to your clipboard
- Add a delay before starting the upload (or before taking the screenshot)

All these features can be enabled/disable with command line flags.

## Configuration file

On the first run, goploader will try to read from 
`~/.config/goploader.conf.yml`. A default configuration will be written in this 
file if it can't be found. Also, if the ~/.config/ directory doesn't exist, it 
will be created. This will be executed only once, you'll get a single line of 
log the first time you'll use the client.

!!! note
    If you want to use another server than mine (a friend of yours installed the
    server or your own server), you'll have to edit this file to point to this 
    service.

## Command line flags

| Short | Long | Default | Description |
|-------|--------------|---------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| -c | --clipboard | false | Copies the returned URL directly to the clipboard (needs xclip or xsel) |
| -d | --delay | 0 | Defines a delay before the program executes. In the form of 10s/ 2m/etc… |
| -n | --name |  | Specify the filename you want. This will be the filename of the file when downloading it or viewing it. When uploading a file, defaults to the name of the file. When uploading from stdin, defaults to stdin. |
| -p | --progress | false | Displays a progress bar while uploading the file. Shows ETA, current network speed, percentage and size uploaded. |
| -s | --screenshot | false | Takes a screenshot and directly uploads it. Works nicely with the delay option. At least one of the following softwares must be installed : `gnome-screenshot`, `import` or `scrot` |
| -t | --tee | false | Displays the content on stdout whever you read from stdin or a file. |
| -v | --verbose | false | Activates the debug logs. Useful for troubleshooting. |
| -w | --window | false | When using the screenshot, this option allows you to select a window or a zone to screenshot before uploading it. |
| -l | --lifetime | 1d | Define the lifetime of your file on the server. Basically you need to choose between 30m, 1h, 6h, 1d, 1w which are the equivalent for 30 minutes, 1 or 6 hours, 1 day, and 1 week. |
| -o | --once | false | Your upload will be visible only once and then deleted from the server |

!!! note "About flags requiring a software"
    Some flags, such as the `-c/--clipboard` require some softwares to be 
    installed on your machine. If it is not avalailable, the client won't work
    with these options, but they are entirely optional.

!!! warning "About Windows and MacOS"
    I was not able to test the screenshot option on Mac or Windows so I don't
    know if they are working properly.

## Examples

This section contains usage examples. Each example is written with both the
long arguments and short arguments. Both are equivalent.

**Take a screenshot of the whole screen after 5 seconds, upload it with a
progress bar and copy the returned URL to the clipboard**

- `$ goploader --progress --screenshot --delay="5s" --clipboard`
- `$ goploader -pscd 5s`

**Select a windows or zone, screenshot that and upload it**

- `$ goploader --screenshot --window`
- `$ goploader -sw`

**Upload a file with a progress bar**

- `$ goploader -p myfile` - File as argument
- `$ goploader -p < myfile` - File as Stdin
- `$ cat myfile | goploader -p` - Pipe

!!! tip "About short arguments"
    Short arguments can be concatenated one after another if and only if they
    don't have to store a value. In short, boolean values can be written one
    after another as shown in the examples.

!!! note "About the progress bar behaviour"
    When the datasource is stdin, the total size of the file (or the total 
    amount of text a command yields) can't be calculated that's why the progress
    bar will only display the upload speed.

 

## File Manager Integration

If you want, you can add an entry to the contextual menu of your file manager to
upload a file. The trick is to use the .desktop files and to copy the returned 
URL to your clipboard once the file is uploaded. Here is an example 
`goploader.desktop` that you can put in your 
`~/.local/share/file-manager/actions/` directory :

```ini
[Desktop Entry]
Type = Action
Tooltip = Upload file on Goploader
Name = Upload on Goploader...
Profiles = goploader_onfile;

[X-Action-Profile goploader_onfile]
MimeTypes = all/allfiles;
SelectionCount = =1
Exec = goploader -c %d/%b
```