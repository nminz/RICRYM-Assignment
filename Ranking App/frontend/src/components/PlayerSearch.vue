<template>
  <div>
    <input 
      v-model="searchTerm" 
      @input="searchPlayers" 
      placeholder="search username" 
      class="search-input"
    />
  </div>
</template>

<script>
export default {
  data() {
    return {
      searchTerm: '' // This will store the current search term
    };
  },
  methods: {
    searchPlayers() {
      // Emit the search term to the parent component
      this.$emit('search', this.searchTerm);
    }
  }
};
</script>

<style scoped>
.search-input {
  background-color: #1a1c1b;
  color: #f0f0f0;
  padding: 10px;
  border: 1px solid #a7955b;
  border-radius: 10dvb;
  width: 20vw;
  min-width: 280px;
  box-sizing: border-box; 
  height: 2rem;
}

@media only screen and (min-width: 1980px) {
  
  .search-input {
  background-color: #1a1c1b;
  color: #f0f0f0;
  padding: 10px;
  border: 1px solid #a7955b;
  border-radius: 10dvb;
  width: 20vw;
  min-width: 280px;
  box-sizing: border-box; 
  height: 4rem;
  font-size: 1.5rem;
}
}
</style>
