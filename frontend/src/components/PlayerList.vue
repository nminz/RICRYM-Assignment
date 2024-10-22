<template>
  <div class="player-list-container">
    <!-- Loading indicator -->
    <div v-if="loading" class="loading">Loading...</div>

    <!-- Table -->
    <table v-else class="player-table" style="max-width: 100%;">
      <thead>
        <tr>
          <th>RANK</th>
          <th @click="sortTable('Username')">
            USERNAME
            <span v-if="sortBy === 'Username'">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
          </th>
          <th @click="sortTable('ClassID')">
            CLASS
            <span v-if="sortBy === 'ClassID'">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
          </th>
          <th @click="sortTable('Score')">
            SCORE
            <span v-if="sortBy === 'Score'">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="player in players" :key="player.AccID">
          <td>{{ player.Rank }}</td>
          <td>{{ player.Username }}</td>
          <td>{{ player.ClassID }}</td>
          <td>{{ player.Score }}</td>
        </tr>
      </tbody>
    </table>

    <!-- Pagination controls -->
    <div class="pagination-controls">
      <button 
        @click="previousPage" 
        :disabled="currentPage === 1"
        class="pagination-button"
      >
        ◄
      </button>
      <span class="page-info" style="margin-bottom: 2rem;">
        {{ currentPage }} / {{ totalPages }}
      </span>
      <button 
        @click="nextPage" 
        :disabled="currentPage === totalPages"
        class="pagination-button"
      >
        ►
      </button>
    </div>
  </div>
</template>

<script>
import { getAccounts } from '@/services/api';

export default {
  data() {
    return {
      players: [],
      currentPage: 1,
      itemsPerPage: 20,
      totalItems: 0,
      totalPages: 0,
      loading: false,
      sortBy: 'Score', ///// default sort = score
      sortOrder: 'desc', ///// default order = descending
      searchTerm: '',
    };
  },

  methods: {
    async fetchPlayers() {
      this.loading = true;
      try {
        const response = await getAccounts(
          this.currentPage,
          this.itemsPerPage,
          this.searchTerm,
          this.sortBy,
          this.sortOrder
        );
        
        this.players = response.data.data;
        this.totalItems = response.data.total;
        this.totalPages = response.data.total_pages;
      } catch (error) {
        console.error('Error fetching players:', error);
      } finally {
        this.loading = false;
      }
    },

    async handleSearch(term) {
      this.searchTerm = term;
      this.currentPage = 1; // Reset to first page when searching
      await this.fetchPlayers();
    },

    async sortTable(column) {
      if (this.sortBy === column) {
        // Toggle sort order if clicking the same column
        this.sortOrder = this.sortOrder === 'asc' ? 'desc' : 'asc';
      } else {
        // New column, set as default desc order
        this.sortBy = column;
        this.sortOrder = 'desc';
      }
      this.currentPage = 1; // Reset to first page when sorting
      await this.fetchPlayers();
    },

    async previousPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
        await this.fetchPlayers();
      }
    },

    async nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++;
        await this.fetchPlayers();
      }
    },
  },

  mounted() {
    this.fetchPlayers();
  },
};
</script>