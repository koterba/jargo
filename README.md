# jargo

1. cross compile using ``GOOS=linux GOARCH=arm64 go build .`` if compiling for the swiftbot/rpi
2. move binary to swiftbot
3. run ``chmod +x jargo`` to make the binary executable
4. run ``sudo mv jargo /usr/local/bin`` to move binary to PATH

## Use

1. ``jargo new <name>`` to create a new project called <name>
2. ``jargo run`` to run the project, once in the directory

if you want to include a java file, drag it into the ``src`` dir.
if you want to include a jar file, drag it into the ``lib`` dir.

## Disclaimer

Main.java and Test.java are placeholder example files in the src directory, same goes for the Example.jar placeholder. I recommend deleting them all except Main.java as you need it as an entrypoint.
