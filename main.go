package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/Aadithya-V/qp-gen/server"
)

//	@title			qp gen service
//	@version		1.0
//	@description	qp gen service
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter the Bearer Authorization string as following: `Bearer Generated-JWT-Token`
//	@scheme						Bearer

func main() {

	// Close db and log file on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// cleanup()
		os.Exit(0)
	}()

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()

	// Start server
	err := server.Start()
	if err != nil {
		panic(err.Error())
	}
	/*
		qs := []Q{
			Q{
				Qq: &Qq{
					Text:  "Sample text 1",
					Marks: 2,
				},
				Id:   1,
				Type: 1,
			},
			Q{
				Qq: &Qq{
					Text:  "Sample text 2",
					Marks: 2,
				},
				Id:   2,
				Type: 1,
			},
			Q{
				Qq: &Qq{
					Text:  "Sample text",
					Marks: 2,
				},
				Id:   2,
				Type: 1,
			},
		}

		QS := BuildQStore(qs)

		qset := QS.PickQSet(map[int]int{
			1: 10,
			2: 2,
			3: 2,
		})

		fmt.Println(qset)

		// Read template from file
		tmpl, err := template.ParseFiles("latex_template.tex")
		if err != nil {
			fmt.Println("Error parsing template:", err)
			return
		}

		// Execute the template with the provided data

		// test transform
		tempQs := make([]string, 0)

		for _, qss := range qset {
			for _, q := range qss {
				tempQs = append(tempQs, q.Qq.Text)
			}
		}
		fmt.Println(tempQs)
		err = tmpl.Execute(os.Stdout, tempQs)
		if err != nil {
			fmt.Println("Error executing template:", err)
			return
		}
	*/
}

type QStore struct {
	ByTypes map[int][]*Q
}

func BuildQStore(Qs []Q) QStore {

	ret := QStore{
		ByTypes: make(map[int][]*Q),
	}

	for _, q := range Qs {
		fmt.Println(q)
		if _, ok := ret.ByTypes[q.Type]; !ok {
			ret.ByTypes[q.Type] = make([]*Q, 0)
		}

		ret.ByTypes[q.Type] = append(ret.ByTypes[q.Type], &q)

	}

	return ret
}

// number of qs per type in opt
func (QS *QStore) PickQSet(opt map[int]int) map[int][]*Q {

	ret := make(map[int][]*Q)

	for qType, qs := range QS.ByTypes {
		if _, ok := ret[qType]; !ok { // if not requrired here
			ret[qType] = make([]*Q, 0)
		}

		ret[qType] = append(ret[qType], pickQ(qs, opt[qType])...)
	}

	return ret
}

func pickQ(qs []*Q, nums int) (ret []*Q) {

	r := rand.Intn(len(qs) - 1)

	for i := r + 1; i < len(qs); i = (i + 1) % len(qs) {
		if nums == 0 {
			return
		}
		if i == r {
			q := qs[rand.Intn(len(qs))]
			ret = append(ret, q)
			q.PickedCount++
			nums--

			continue
		}

		if qs[i].PickedCount == 0 {
			ret = append(ret, qs[i])
			qs[i].PickedCount++
			nums--
		}

	}
	return
}
