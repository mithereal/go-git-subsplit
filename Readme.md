Go Git-Subsplit

by Jason Clark (mithereal@gmail.com)

Automate and simplify the process of managing one-way read-only subtree splits concurrently via go.

Git subsplit relies on subtree being available. If is not available in your version of git (likely true for versions older than 1.7.11) please install it manually.


Can sync Branches, Tags and Origins.

This script was inspired by the following

https://github.com/dflydev/git-subsplit

https://github.com/cebe/git-simple-subsplit

reason for creating: these versions already existed that mainly work, except when dealing with large repositeries, with go we can 
manage many concurrent processes in order to exponentially increase productivity.
