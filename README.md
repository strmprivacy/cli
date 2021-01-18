# Stream Machine Command Line Interface

## Installation

1. Make sure that you put the `strm` binary somewhere on your `PATH`.
2. (only for initial setup) Add the following lines to your `.bashrc` / `.zshrc` or equivalent: `source <(strm --generate-completion <shell>)`. Don't forget to replace `<shell>` with the respective value for your shell type (`bash`, `zsh`, `fish`).

## Developer setup

1. Install GraalVM: `sdk list java` and `sdk install java <latest_graalvm>`. Note: use `r11` for JDK11+ and `r8` for JDK8
2. Install Native Image GraalVM component: `gu install native-image`
3. Make sure when running maven, you've set GraalVM as the JDK to use in that shell (hence, `JAVA_HOME` is set to GraalVM)

