## Quick Doc

Simply `curl https://gplr.in` to see a small reminder/documentation on how to
use the service. 

## Usage

One of the main goal of goploader is the ease of use. That's also why the 
goploader server doesn't use a JSON API, to guarantee that any user can upload 
files painlessly.

You can use this service without a client, simply by issuing curl commands.

**Upload the file named `myfile.txt`**

`$ curl -F file=@myfile.txt https://gpldr.in/`

**Change the name of the file to "myamazingfile!"**

`$ curl -F name="myamazingfile!" -F file=@myfile.txt https://gpldr.in/`

**Upload from stdin**

`$ tree | curl -F file=@- https://gpldr.in/`

**Upload a file from stdin**

`$ curl -F file=@- https://gpldr.in/ < myfile.txt`

**Specify the lifetime of your file**

`$ curl -F file=@- -F duration=1w https://gpldr.in/ < myfile.txt`

**Specify that the file is visible only once**

`$ curl -F file=@- -F once="true" https://gpldr.in/ < myfile.txt`

!!! note "Lifetime options"
    See the [Client/Documentation](client/documentation.md) for more information
    about the possible values of the `duration` field.

## File Manager Integration

If you want, you can add an entry to the contextual menu of your file manager 
to upload a file. The trick is to use the .desktop files and to copy the 
returned URL to your clipboard once the file is uploaded. Here is an example 
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
Exec = curl -F file=@%d/%b https://gpldr.in/ | xclip
```

!!! note 
    You will need xclip to copy the returned URL to your clipboard.