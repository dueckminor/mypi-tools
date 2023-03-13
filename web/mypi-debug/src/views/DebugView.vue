<template lang="pug">
v-container
  v-expansion-panels
    v-expansion-panel(:title="service.name" v-for="service in services")
      v-expansion-panel-text
        v-expansion-panels(multiple)
          v-expansion-panel(v-for="component in service.components")
            v-expansion-panel-title
              div
                v-chip(v-if="component.state=='running'",color="green") {{ component.state }}
                v-chip(v-else-if="component.state=='stopped'",color="red") {{ component.state }}
                v-chip(v-else,color="orange") {{ component.state }}
              div &nbsp; {{ component.name }}
              v-spacer
              template(v-for="action in component.actions")
                v-btn(v-if="action.name=='restart'" @click="call_action(service.name,component.name,action.name)",@click.native.stop)
                  v-icon mdi-restart
                v-btn(v-else-if="action.name=='debug'" @click="call_action(service.name,component.name,action.name)",@click.native.stop)
                  v-icon mdi-bug-outline
                v-btn(v-else @click="call_action(service.name,component.name,action.name)",@click.native.stop) {{action.name}}
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
function call_action(service,component,action) {
  // eslint-disable-next-line no-console
  console.log("call_action "+component+action);
  axios({ method: "POST", url: "/api/services/"+service+"/components/"+component+"/"+action })
}
</script>
<!-- ----------------------------------------------------------------------- -->
<!--
<style scoped>
</style>
-->