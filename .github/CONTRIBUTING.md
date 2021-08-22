# How To Contribute

This is an open source project, run by volunteers in their free time. To keep this environment as 
pleasant as possible for everyone, please adhere to the [Go code of conduct](https://golang.org/conduct).

## Finding something to do

If you know what you want to do, feel free to skip this.

If you're looking for something, have a look at the 
[list of raised issues](https://github.com/PaulSonOfLars/gotgbot/issues) - This might help provide some inspiration.

You can also contribute by improving documentation, which will improve the docs available on https://go.pkg.dev.

## Making a code change
- Make sure you've checked out the correct branch to make changes to. The latest version is currently the `v2` branch. 
  This is also the version you will open the PR against.
- Run the `./scripts/setup-hooks.sh` script to ensure your code gets linted before pushing to github. This will run a 
  set of go linters to apply various codestyle guidelines, which will save time during PR review.
- Make sure new functions are documented!
- Please make sure to fill out the PR template to describe your changes. This makes it easier for others to review, 
  and will help get your PR merged sooner.
