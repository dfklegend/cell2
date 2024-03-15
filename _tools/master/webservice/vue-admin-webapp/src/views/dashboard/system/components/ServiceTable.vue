<template>
  <div class="className">
    <el-card class="anoCard">
      <div class="searchDiv"></div>
      <el-table :data="tableData" border stripe>
        <!--      <el-table-column type="selection" width="55"></el-table-column>-->
        <el-table-column label="节点" prop="node" width="140"></el-table-column>
        <el-table-column
          label="运行时间（秒）"
          prop="time"
          width="140"
        ></el-table-column>
        <el-table-column label="服务" prop="service"></el-table-column>
        <el-table-column
          label="状态"
          prop="status"
          width="100"
        ></el-table-column>
        <el-table-column label="操作" width="300">
          <template #default="scope">
            <el-button type="primary" @click="onClickRetire(scope.row.node)"
              >退休
            </el-button>
            <el-button type="danger" @click="onClickExit(scope.row.node)"
              >退出
            </el-button>
            <el-button type="primary" @click="loadNodeData(scope.row)"
              >节点详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-dialog title="节点详情" :visible.sync="diaIsShow" class="diaForm">
      <el-table :data="nodeData" border stripe>
        <el-table-column prop="_key" label="指标"></el-table-column>
        <el-table-column prop="_value" label="数据"></el-table-column>
      </el-table>
      <br />
      <el-button type="danger" @click="diaIsShow = false">关闭</el-button>
    </el-dialog>
  </div>
</template>

<script>
import $axios from '@/api/index'

export default {
  data() {
    return {
      tableData: [],
      status: '',
      diaIsShow: false,
      nodeData: []
    }
  },

  mounted() {
    this.getAllNodes()
  },

  methods: {
    onClickRetire(node) {
      this.retireNode(node)
    },
    onClickExit(node) {
      this.exitNode(node)
    },
    getAllNodes() {
      $axios
        .get('/api/nodes')
        .then(res => {
          this.tableData = res
        })
        .catch(error => {
          this.$message.error(error.message)
        })
    },
    retireNode(node) {
      $axios
        .post('/api/retire', { p: node })
        .then(res => {
          this.$message.success(res)
        })
        .catch(error => {
          this.$message.error(error.message)
        })
    },
    exitNode(node) {
      $axios
        .post('/api/exit', { p: node })
        .then(res => {
          this.$message.success(res)
        })
        .catch(error => {
          this.$message.error(error.message)
        })
    },
    loadNodeData(data) {
      this.diaIsShow = true
      if (data.nodeData) {
        let map = data.nodeData
        this.nodeData = []
        for (const dataKey in map) {
          this.nodeData.push({
            _key: dataKey,
            _value: map[dataKey]
          })
        }
      } else {
        this.nodeData = []
      }
    }
  }
}
</script>
<style lang="scss" scoped>
.fyDiv {
  float: right;
  margin-top: 30px;
  padding-bottom: 20px;
}

.anoCard .el-table .el-button {
  padding: 8px 18px;
  font-size: 12px;
}

.searchDiv {
  margin-bottom: 20px;

  .el-button {
    padding: 11px 20px;
  }
}

.width1 {
  width: 180px;
  margin-right: 10px;
}

.diaForm {
  .el-input {
    width: 350px;
  }
}
</style>
<style lang="scss">
.anoCard {
  .el-card__body:after {
    content: '';
    clear: both;
    width: 0;
    height: 0;
    visibility: hidden;
    display: block;
  }
}

.diaForm .el-form-item__label {
  padding-right: 20px;
}

.searchDiv [class^='el-icon'] {
  color: #fff;
}
</style>
