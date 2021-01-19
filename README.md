# Stream Machine Command Line Interface

## Installation

(commands for Unix platforms included)

1. Download the binary for your platform from the latest release (https://github.com/streammachineio/cli/releases).
2. Rename the binary to `strm` and make sure it's executable (`chmod +x strm`)
3. Make sure that you put the `strm` binary somewhere on your `PATH`. (e.g. `mv strm /usr/local/bin`)
4. (only for initial setup) Add the following lines to your `.bashrc` / `.zshrc` or equivalent: `source <(strm --generate-completion <shell>)`. Don't forget to replace `<shell>` with the respective value for your shell type (`bash`, `zsh`, `fish`).
5. For macOS: please allow the binary to circumvent app notarization for now (we're looking into other solutions): `xattr -d -r com.apple.quarantine strm mac`

## Developer setup

1. Install GraalVM: `sdk list java` and `sdk install java <latest_graalvm>`. Note: use `r11` for JDK11+ and `r8` for JDK8
2. Install Native Image GraalVM component: `gu install native-image`
3. Make sure when running maven, you've set GraalVM as the JDK to use in that shell (hence, `JAVA_HOME` is set to GraalVM)
