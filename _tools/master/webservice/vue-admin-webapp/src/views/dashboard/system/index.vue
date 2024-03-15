<template>
  <div>
    <h1 style="text-align: center">Master 节点服务 监控</h1>
    <br />
    <service-table></service-table>
  </div>
</template>

<script>
import ServiceTable from './components/ServiceTable'

export default {
  components: {
    ServiceTable: ServiceTable
  },

  data() {
    return {
      status: ''
    }
  },

  mounted() {},

  methods: {}
}
</script>

<style scoped></style>
