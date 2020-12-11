import Vue from 'vue'
/* tslint:disable */
<% options.components.forEach(({ name, path }) => { %>
  import <%= name %> from '<%= path %>'
<% }) %>

<% options.components.forEach(({ name }) => { %>
  Vue.component('<%= name %>', <%= name %>)
<% }) %>

// console.log(`[CB Components] Plugin called...`)
/* tslint:enable */
