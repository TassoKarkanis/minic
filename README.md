# MiniC

A toy AMD64 C compiler implemented in Go.

## Developing with Visual Studio Code

#### Step 1. Install [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension

This extension is required to enable dev containers in your native VS Code.

#### Step 2. Build and start the dev container

Run the command:
```
make devcontainer
```
The dev container (named "minic-devcontainer") will be built and started.

#### Step 3. Connect to the dev container

In VS Code, press F1 or Ctrl-Shift-P to open the command palette and run the command "Dev Containers: Attach to Running Container...".  Choose the container "/minic-devcontainer".  A new VS Code window will open.  In this window, open the folder "/w" (this is where your clone of the repo is mounted).

#### Step 4. Configure the dev container

In the new VS Code window that opened in Step 3, run the command "Dev Containers: Open Attached Container Configuration File...".  Choose the container "/minic-devcontainer".  Copy and paste this into the configuration file:
```
{
  "workspaceFolder": "/w",
  "extensions": [
    "golang.Go",
  ],
  "remoteUser": "root"
}
```

#### Step 5. Reset the dev container

Close the VS Code window that opened in Step 3.  Re-run the command from Step 2.  Attach to the container again as in Step 3.  The VS Code window should open to the Acorn Platform project and the extensions listed above should be installed.


## Unimplemented Compiler Features

* Logical operators
  * comparison operators: == and !=
  * binary operators: && and ||
  * unary operators: !
* Basic pointers
  * pointer variables
  * address-of operator: &
  * dereference operator: *
  * offset/indexing with []
* Functions
  * function calls
  * function declarations
  * vararg function calls
* Bitwise operators
  * binary operators: & | ^
  * unary operators: ~
* Other integer types
  * char
  * unsigned char
  * short
  * unsigned short
  * unsigned int
  * long
  * unsigned long
  * long long
  * unsigned long long
* Floating point
  * addition
  * subtraction
  * unary minus
  * multiplication
  * division
* Strings
  * char literals
  * string literals
* Flow control
  * if/else statements
  * while loops
  * for loops
  * do/while loops
  * break statements in loops
  * continue statements
  * switch statements
  * goto statements
* Literals
  * char literals
  * char literals with escapes
  * floating decimal literals
  * floating point exponent literals
  * integer hex literals
  * integer octal literals
  * struct literals
* Global variables
  * declarations
  * definitions without initializers
  * initializers
* Structs
  * declarations
  * basic type field evaluation
  * basic type field field assignment
  * struct rvalue/lvalue
* Arrays
  * declarations
  * definitions
  * initializers
* Other
  * sizeof operator
  * type casts
  * typedef
  * pointer to function
  * ternary operator
  * bit fields
  * static functions
  * static global variables