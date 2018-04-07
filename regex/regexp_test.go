package regex

import (
	"testing"
	"regexp"
	"fmt"
	"bytes"
)

//formula is define for the name regex.
var validNameRegex = regexp.MustCompile(`^[\.a-zA-Z0-9_-]+(#ephemeral)?$`)


func TestSimpleRegExp(t *testing.T){
	//This case is often used.
	fmt.Println("'topic-name' is valid name", validNameRegex.MatchString("topic-name"))
	//This tests whether a pattern matches a string.
	match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
	fmt.Println("p([a-z]+)ch, peach, match:", match)


	r, _ := regexp.Compile("p([a-z]+)ch")
	//This finds the match for the regexp.
	fmt.Println("find match: ", r.FindString("peach punch peach"))

	//We can also provide []byte arguments and drop String from the function name.
	fmt.Println(r.Match([]byte("peach")))

	//This also finds the first match but returns the start and end indexes for the match instead of the matching text.
	fmt.Println("find match: ", r.FindStringIndex("peach punch peach"))

	//The Submatch variants include information about both the whole-pattern matches and the submatches within those matches. For example this will return information for both p([a-z]+)ch and ([a-z]+).
	fmt.Println("find match: ", r.FindStringSubmatch("peach punch peach"))
	fmt.Println("find match: ", r.FindStringSubmatchIndex("peach punch peach"))

	//These All variants are available for the other functions we saw above as well.
	fmt.Println(r.FindAllString("peach good punch idiom patch pinch", -1))

	//The regexp package can also be used to replace subsets of strings with other values.
	fmt.Println("peach is replaced by fruit", r.ReplaceAllString("a peach", "fruit"))

	//The Func variant allows you to transform matched text with a given function.
	in := []byte("a peach")
	fmt.Println("replace by func , ", string(r.ReplaceAllFunc(in, bytes.ToUpper)))
}

func TestAdvancedRegExp(t *testing.T){
	//[[cat c] [sat s] [mat m]]
	re, err := regexp.Compile(`(.)at`) // want to know what is in front of 'at'
	if err != nil {
		t.Fatal(err)
	}
	res := re.FindAllStringSubmatch("The cat sat on the mat.", -1)
	fmt.Printf("%v", res)
	//match case
	fmt.Println("\n", "match case: ", re.FindAllString("The cat sat on the mat.", -1))

	//match http domain name
	dmRegexp, _ := regexp.Compile(`([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}`) // want to know what is domain
	domain := dmRegexp.FindString("https://research.swtch.com/go2017")
	if domain != "research.swtch.com" {
		t.Fatal("error return domain: ", domain)
	}
	oftenDomain := dmRegexp.FindString("http://dict.cn/often")
	if oftenDomain != "dict.cn" {
		t.Fatal("error return domain: ", domain)
	}

	//Literal Special Characters
	//Finding one backslash '\': It must be escaped twice in the regex and once in the string.
	backslash, _ := regexp.Compile("C:\\\\")
	if backslash.MatchString("working on drive C:\\") == true {
		fmt.Printf("Matches.") // <---
	}else {
		fmt.Printf("No match.")
	}
	//The other special characters that are relevant for constructing regular expressions work in a similar fashion: .+*?()|[]{}^$

	//Simple Repetition
	s := "Firstname Lastname"
	r, err := regexp.Compile(`\w+\s\w+`)
	// Prints Firstname Lastname
	fmt.Printf("%v \n",  r.FindString(s))
}

