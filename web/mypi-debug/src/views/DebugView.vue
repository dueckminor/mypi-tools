<template lang="pug">
v-container
  v-expansion-panels
    v-expansion-panel(:title="service.name" v-for="service in services")
      v-expansion-panel-text
        v-expansion-panels(multiple)
          v-expansion-panel(v-for="component in service.components")
            v-expansion-panel-title
              v-chip(v-if="component.state=='running'",color="green") {{ component.state }}
              v-chip(v-else-if="component.state=='stopped'",color="red") {{ component.state }}
              v-chip(v-else,color="orange") {{ component.state }}
              v-btn(@click="restart_component(service.name,component.name)",@click.native.stop) restart
              div - {{ component.name }}
            v-expansion-panel-text
              go-tty(:path="`/api/services/`+service.name+`/components/`+component.name+`/tty`" style="padding: 0")
</template>
<!-- ----------------------------------------------------------------------- -->

<script setup>
import { storeToRefs } from 'pinia'
import axios from "axios";
// eslint-disable-next-line no-unused-vars
import GoTty from "../components/GoTTY";
import { useServiceStore } from "@/stores/ServiceStore";

// eslint-disable-next-line no-unused-vars
const services = storeToRefs(useServiceStore()).services
//const services = useServiceStore()

// eslint-disable-next-line no-unused-vars
function restart_component(service,component) {
  // eslint-disable-next-line no-console
  console.log("restart_component "+component);
  axios({ method: "POST", url: "/api/services/"+service+"/components/"+component+"/restart" })
}
</script>
<!-- ----------------------------------------------------------------------- -->
<!--
<style scoped>
</style>
-->