
### string到int  
  int,err:=strconv.Atoi(string)  
### int到string  
 string:=strconv.Itoa(int)  
 
### string到int64  
  int64, err := strconv.ParseInt(string, 10, 64)  
### int64到string  
 string:=strconv.FormatInt(int64,10)  
 
 
 
 
### Install the RDM
```
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)" < /dev/null 2> /dev/null ; brew install caskroom/cask/brew-cask 2> /dev/null


brew cask install rdm
```
