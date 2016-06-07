package gitremote

/*
 * The aim of the package is simply to parse the remotes of a .git/config
 * file. To get the differents names and their destination
 *
 * remotes, err := gitremote.List()
 * ---
 *
 * Config.Add will return an error if the remote already exists
 *
 * config := giremote.New("repo/path")
 * err := config.Add(&gitremote.Remote{
 *   Name: "upstream",
 *   URL: "git@github.com:Soulou/go-gitremote.git"
 * })
 * ---
 *
 * Config.AddOrUpdate is like Add, but will update an existing remote
 *
 * config := giremote.New("repo/path")
 * err := config.AddOrUpdate(&gitremote.Remote{
 *   Name: "upstream",
 *   URL: "git@github.com:Soulou/go-gitremote.git"
 * })
 * ---
 *
 * Config.Delete will remote a remote from a git config
 *
 * config := giremote.New("repo/path")
 * err := config.Delete(&gitremote.Remote{Name: "origin"})
 */
