package main

const BLAISE = "[^*BLAISE*^]"

func BroadcastUpdate(packageName string) {
	for packageName, connection := range Clients {
		_, err := connection.Write([]byte(BLAISE+"\n"))
			if err != nil {
				// TODO: add to queue, to deliver eventually
				connection.Close()
				delete(Clients, packageName)
			}
		}
}