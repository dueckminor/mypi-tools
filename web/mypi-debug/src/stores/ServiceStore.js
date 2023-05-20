import { defineStore } from "pinia";
import axios from "axios";
//import io from "socket.io-client";


export const useServiceStore = defineStore("ServiceStore", {
    // state
    state: () =>{
        return {
            error: null,
            services:[]
        }
    },
    // actions
    actions: {
        async fill() {

            axios({ method: "GET", url: "/api/services" }).then(
                result => {
                    this.services = result.data
                    const socket = new WebSocket('ws://localhost:8080/api/ws/')
                    socket.onmessage = (event) => {
                        const msg = JSON.parse(event.data)
                        const parts = msg.topic.split("/")
                        this.services.forEach(service => {
                            if (service.name == parts[0]) {
                                service.components.forEach(component => {
                                    if (component.name == parts[1]) {
                                        component[parts[2]]=msg.value
                                    }
                                })
                            }
                        })
                    }
                },
                error => {
                    this.error = error
                }
              );
        }
    }
    // getters
} )