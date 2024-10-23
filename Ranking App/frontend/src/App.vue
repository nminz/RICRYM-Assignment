<template>
  <div id="app">
    <img src="./assets/wira-logo.png" alt="WIRA" style="width:10%; min-width: 120px; margin-top: 3rem; justify-content: center; align-items: center;">
    <h1>GLOBAL RANKING</h1>
    <img src="./assets/text-underline.png" alt="underline" style="width: 20%; min-width: 300px; margin: 1rem;">


    <PlayerSearch @search="handleSearch" />
    <PlayerList ref="playerList" />

  </div>
</template>

<script>

import PlayerList from './components/PlayerList.vue';
import PlayerSearch from './components/PlayerSearch.vue';
import '@/assets/css/main.css'

export default {
  name: 'App',
  components: {
    PlayerList,
    PlayerSearch
  },
  methods: {
    handleSearch(searchTerm) {

      this.$refs.playerList.handleSearch(searchTerm);
    }
  }
}
</script>