<template>
    <v-container fluid fill-height>
        <GoTTY path="/ws/terminal"></GoTTY>
    </v-container>
</template>


<script>

import GoTTY from '../components/GoTTY'

export default {
    name: 'initialize_sd',
    components: {
        GoTTY
    },
    computed: {
    },
    mounted() {
    },
    methods: {
    }
};
</script>

<style scoped></style>
