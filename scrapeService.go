package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type ScrapeSvc struct {
	DbCtrl       *DbCtrl
	statusMap    map[string]bool
	mdataChannel chan string
}

func NewScrapeSvc() *ScrapeSvc {
	return &ScrapeSvc{DbCtrl: NewDbCtrl(),
		statusMap: make(map[string]bool), mdataChannel: make(chan string)}
}

func (p *ScrapeSvc) fetchMeta() {

}

func (p *ScrapeSvc) Run(cids []string, done chan struct{}, s *SharedExtConn) error {

	log.Println("Scraping Service Started!")

	//Read through CID array and fetch
	for i := 0; i < len(cids); i++ {
		p.mdataChannel <- cids[i]
	}

	concurrent_threads := 3
	//Write to channel
	for i := 0; i < concurrent_threads; i++ {
		go p.fetchMeta()
	}

processingLoop:
	for {
		select {
		//process incoming messages from broker

		//app provision request to backend
		//this should go to AppManagerSvc
		//metadata: get an address for account (pki), funding amount for account
		case t4 := <-ctrlChannel:

			if t4.msgType == "APPPROV" {

			}

		//coming from web-tier
		case t3 := <-ctrlChannelPayReq:

			if t3.msgType == "PAYREQ" { //Transaction request

				pr := t3.body
				log.Printf("Read payment id: %s\n", pr.TxId)
				fmt.Printf("Read payment id: %s\n", pr.TxId)

				intVar, err := strconv.Atoi(pr.NodeIdFrom)

				if err == nil {
					acctinfo, err := s.Msql.GetAccountForDevice(int64(intVar))
					if err == nil {
						out, err := p.LndCtrl.sendPaymentCustodial(acctinfo.accountName, acctinfo.host, acctinfo.port, pr)
						if err == nil {
							//parse out
							inv := lPayment{}
							json.Unmarshal([]byte(out), &inv)
							lnc <- bcNotification{msgType: "LNDPaymentInitiated", body: pr.TxId}
							log.Println("LNDPaymentInitiated")
							// DB Operations
						}
					}
				}
			}

		// { "balances": [{"nodeId": "", "balance": ""},{"nodeId": "", "balance": ""}]
		// move this to webservice
		/*case t2 := <-ctrlChannelr:

		if t2.msgType == "GWBR" { //wallet balance

			//ctrlChannel <- ctrlNotification{msgType: "GWB", body: []WalletBalance{{"1", p.dvcWallet["1"]},
			//	{"2", p.dvcWallet["2"]}}}

			//p.dvcWallet["1"] = 1000
			//p.dvcWallet["2"] = 2000

			balanceOut := `{"balances":[{"nodeId":"1","balance":"` + strconv.Itoa(p.dvcWallet["1"]) + `"},` + `{"nodeId":"2","balance":"` + strconv.Itoa(p.dvcWallet["2"]) + `"}]}`
			//ctrlChannel <- ctrlNotification{msgType: "GWB", body: "{  \"balances\": [{\"nodeId\": \"1\", \"balance\": p.dvcWallet[\"1\"]}] }"}
			//ctrlChannel <- ctrlNotification{"GWB", `{"balances":[{"nodeId":"1","balance":"1000"},{"nodeId":"2","balance":"10000"}]}`}
			ctrlChannel <- ctrlNotification{"GWB", balanceOut}

		}
		*/

		case t := <-dvcChannel:
			//Read the message from device
			if t.msgType == "DAuth" {
				app, err := s.Msql.GetApp(t.body)
				fmt.Printf("Received app id: %s\n", app)
				if err == nil {
					intVar, _ := strconv.Atoi(app)
					if intVar > 0 {
						text := fmt.Sprintf("dregis")
						token := p.mqttClient.Publish("lnd/pctrl", 0, false, text)
						log.Println("DRegis")
						token.Wait()
						//update db with newly registered device
						// map to store flattened struct
						var d = make(map[string]string)
						// flatten the GET params
						d["dappid"] = app
						d["ddesc"] = " "
						d["downer"] = " "

						l := time.Now()
						b := float64(0)
						st := true

						// build the device struct
						device := &Device{W: nil, R: nil, D: d, S: s, LastActive: l, Status: st, Balance: b}
						status, err := processDeviceRegistrationMqtt(device, l, b, st)
						if err == nil {
							//log.Printf("ERROR: Processing Failed, Returning (%s)", "http.StatusBadRequest")
							//http.Error(w, err.Error(), http.StatusBadRequest)
							//return

							//Create invoice and push to db
							intVar, err := strconv.Atoi(status)
							log.Printf("Read device id: %d\n", intVar)
							fmt.Printf("Read device id: %d\n", intVar)
							if err == nil {
								acctinfo, err := s.Msql.GetAccountForDevice(int64(intVar))
								log.Printf("Account name: %v\n", acctinfo.accountName)
								if err == nil {
									fmt.Println("Acct name error is nil")
									out, err := p.LndCtrl.addInvoiceCustodial(acctinfo.accountName, acctinfo.host, acctinfo.port)
									if err == nil {
										log.Println("LNDInvoiceCreated ready")
										fmt.Printf("outbefore:\n%s\n", out)
										//parse out
										inv := lInvoice{}
										//json.Unmarshal([]byte(out), &inv)
										//log.Printf("r_hash = %v\n", inv.rHash)
										//fmt.Printf("r_hash = %v\n", inv.rHash)

										if err := json.Unmarshal([]byte(out), &inv); err != nil {
											fmt.Println("failed to unmarshal:", err)
										} else {
											log.Printf("r_hash = %v\n", inv.RHash)
											fmt.Printf("payment_request = %v\n", inv.PaymentRequest)
										}

										//1 is a dummy device id
										/*
											log.Println("Here1")
											p.dvcMap["1"] = inv.PaymentRequest
											log.Println("Here2")
											lnc <- bcNotification{msgType: "LNDInvoiceCreated", body: "1"}
										*/

										//store it in db
										id, err := s.Msql.InsertInvoiceForDevice(int64(intVar), inv.PaymentRequest)
										if err == nil {
											log.Printf("Invoice inserted with id: %d", id)
										}

										log.Println("LNDInvoiceCreated")
									}
								}
							}
						}

						//check this
						//w.Header().Set("Content-Type", "application/json")

						// write back to the client
						//_, err = fmt.Fprintf(w, "%s", status)
						//log.Printf("Status [%s]", status)

					}
				}
			}

			if t.msgType == "DInvoice" {

				//store current wallet value of Charlie
				/*
					out1, err1 := p.LndCtrl.listChannelsN()
					if err1 == nil {
						lchs := lChannels{}
						json.Unmarshal([]byte(out1), &lchs)
						currBalance, _ := strconv.Atoi(lchs.Channels[0].LocalBalance)
						p.dvcWallet["1"] = currBalance
						currBalanceR, _ := strconv.Atoi(lchs.Channels[0].RemoteBalance)
						p.dvcWallet["2"] = currBalanceR
					}*/

				//create invoice
				//write to lnc
				log.Println("Read device message")
				fmt.Println("Read device message101")
				//out, err := p.LndCtrl.addInvoice()

				/* Query DB to find out custodial account
				------------------------------------------
				*/
				intVar, err := strconv.Atoi(t.body)
				log.Printf("Read device id: %d\n", intVar)
				fmt.Printf("Read device id: %d\n", intVar)
				if err == nil {
					acctinfo, err := s.Msql.GetAccountForDevice(int64(intVar))
					log.Printf("Account name: %v\n", acctinfo.accountName)
					if err == nil {
						fmt.Println("Acct name error is nil")
						out, err := p.LndCtrl.addInvoiceCustodial(acctinfo.accountName, acctinfo.host, acctinfo.port)
						if err == nil {
							log.Println("LNDInvoiceCreated ready")
							fmt.Printf("outbefore:\n%s\n", out)
							//parse out
							inv := lInvoice{}
							//json.Unmarshal([]byte(out), &inv)
							//log.Printf("r_hash = %v\n", inv.rHash)
							//fmt.Printf("r_hash = %v\n", inv.rHash)

							if err := json.Unmarshal([]byte(out), &inv); err != nil {
								fmt.Println("failed to unmarshal:", err)
							} else {
								log.Printf("r_hash = %v\n", inv.RHash)
								fmt.Printf("payment_request = %v\n", inv.PaymentRequest)
							}

							//1 is a dummy device id
							/*
								log.Println("Here1")
								p.dvcMap["1"] = inv.PaymentRequest
								log.Println("Here2")
								lnc <- bcNotification{msgType: "LNDInvoiceCreated", body: "1"}
							*/

							//store it in db
							id, err := s.Msql.InsertInvoiceForDevice(int64(intVar), inv.PaymentRequest)
							if err == nil {
								log.Printf("Invoice inserted with id: %d", id)
							}

							log.Println("LNDInvoiceCreated")
						}
					}
				}

			}
		//process incoming messages from lncclient1
		//Read from lnc about valid payments on bc
		case t1 := <-lnc:

			//client.Publish
			/*if t1.msgType == "LNDInvoiceCreated" {
				//parse out
				invId := p.dvcMap["1"]

				// send payment
				// move this to API Server ?
				out, err := p.LndCtrl.sendPayment(invId) // ??needs wallets
				if err == nil {
					//parse out
					inv := lPayment{}
					json.Unmarshal([]byte(out), &inv)
					lnc <- bcNotification{msgType: "LNDPaymentInitiated"}
					log.Println("LNDPaymentInitiated")
				}
			}*/

			//check BTC ledger every few seconds
			//How does custodial account map wallet balance to a device
			//bUG FIX?
			if t1.msgType == "LNDPaymentInitiated" {
				acctinfo, err := s.Msql.GetAccountForTx(t1.body)
				log.Printf("Account name: %v\n", acctinfo.accountName)
				if err == nil {
					//check ledger balance
					out, err := p.LndCtrl.listChannelsCustodial(acctinfo.accountName, acctinfo.host, acctinfo.port)
					if err == nil {
						//parse out
						//lchs := make([]lChannel, 0)
						lchs := lChannels{}
						json.Unmarshal([]byte(out), &lchs)
						log.Printf(" lchs.Channels[0].LocalBalance = %v\n", lchs.Channels[0].LocalBalance)
						//Need to fix this ?
						i, _ := strconv.ParseFloat(lchs.Channels[0].LocalBalance, 8)
						//iRemote, _ := strconv.ParseFloat(lchs.Channels[0].RemoteBalance, 8)
						log.Printf(" lchs.Channels[0].Remotebalance = %v\n", lchs.Channels[0].RemoteBalance)

						//if i == 900950 /*CURR_BALANCE-PAY_AMT)*/ {
						// if i == CURR_BALANCE-PAY_AMT {

						invoiceAmt, err := s.Msql.GetTransactionAmountForTx(t1.body)
						if err == nil {
							amt, _ := strconv.Atoi(invoiceAmt)
							currBalance, _ := s.Msql.GetBalanceForPayee(t1.body)
							log.Printf("Current balance of payee = %v\n", currBalance)
							if i == float64(amt)+currBalance {
								//if iRemote == float64(amt)+currBalance {
								//update wallets
								_, _ = s.Msql.UpdateWallets(t1.body)
								lnc <- bcNotification{msgType: "LNDPaymentCommitted"}
								log.Println("LNDPaymentCommitted")

							} else {
								time.Sleep(2 * time.Second)
								lnc <- bcNotification{msgType: "LNDPaymentInitiated", body: t1.body}
								log.Println("LNDPaymentInitiated again")
							}
						}
						//Compare if balance incremented
						/*if i == p.dvcWallet["1"]+PAY_AMT {
							lnc <- bcNotification{msgType: "LNDPaymentCommitted"}
							log.Println("LNDPaymentCommitted")
							//update wallet balances
							p.dvcWallet["1"] = i
							p.dvcWallet["2"] = -PAY_AMT

						} else {
							time.Sleep(2 * time.Second)
							lnc <- bcNotification{msgType: "LNDPaymentInitiated"}
							log.Println("LNDPaymentInitiated again")
						}*/
					}

				}
			}

			if t1.msgType == "LNDPaymentCommitted" {

				text := fmt.Sprintf("dp")
				//token := p.mqttClient.Publish("topic/test1", 0, false, text)
				token := p.mqttClient.Publish("lnd/pdone", 0, false, text)
				log.Println("DPayment")
				token.Wait()
				//write to control channel also for web tier response to client
				ctrlChannel <- ctrlNotification{"PAYRES", "Transaction Complete !"}
				//check done
				//done <- struct{}{}

			}

		case <-done:

			break processingLoop

		}
	}

	return nil

}
