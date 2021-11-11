# @customrealms/cli

CustomRealms command-line tools for setting up and compiling JavaScript Minecraft plugins.

### Installation

Install the CLI on your computer:

```sh
npm install -g @customrealms/cli
```

Then, you can use the CustomRealms CLI using the `crx` command in your terminal.

### Start a project

```sh
mkdir my-plugin
cd my-plugin
crx init
```

That's it! You now have a plugin project ready to develop!

### Build a JAR file

Compile your plugin project to a JAR file:

If you used `crx init` above to create your project, you can compile a JAR file using:

```sh
npm run build:jar
```

If you used a different method to create your project, use this instead:

```sh
crx build -o ./dist/my-plugin.jar
```
