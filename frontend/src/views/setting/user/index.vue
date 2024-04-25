<template>
  <div>
    <el-card shadow="never">
      <div slot="header" class="clearfix">
        用户/角色列表
      </div>
      <div class="list">
        <el-row style="margin-bottom: 20px">
          <el-col>
            <el-button
              icon="el-icon-plus"
              type="primary"
              size="small"
              @click="createServer"
              >增加授权
            </el-button>
          </el-col>
        </el-row>
        <el-table :data="tableData" height="600" border stripe>
          <el-table-column prop="user" label="用户/角色"></el-table-column>
          <el-table-column label="权限">
            <template slot-scope="scope">
              <el-tag v-for="(item, idx) in scope.row.roles" :key="idx">{{
                item
              }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="是否启用">
            <template slot-scope="scope">
              <el-switch v-model="scope.row.is_available" @change="handleAble(scope.row)">
              </el-switch>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-button
                size="mini"
                type="primary"
                icon="el-icon-eidt"
                @click="edit(scope.row)"
                >编辑</el-button
              >
            </template>
          </el-table-column>
        </el-table>

        <el-pagination
          background
          :page-sizes="[5, 10, 50, 100]"
          :page-size="10"
          layout="total, sizes, prev, pager, next, jumper"
          :total="0"
        >
        </el-pagination>
      </div>
    </el-card>

    <FormModel ref="form" @handleLoading="handleLoading" />
  </div>
</template>

<script>
import FormModel from "./components/form";

export default {
  name: "User",
  components: {
    FormModel
  },
  data() {
    return {
      tableData: [],
      isAble: false
    };
  },
  created() {
    this.getUserRole();
  },
  methods: {
    createServer() {
      this.$refs.form.show();
    },
    handleLoading(data) {
      if (!data) {
        this.getUserRole();
      }
    },
    async getUserRole() {
      try {
        let res = await this.$api.user.GetAllUser();
        res.forEach(item => {
          if (item.del === "0") {
            this.$set(item, "is_available", true);
          } else {
            this.$set(item, "is_available", false);
          }
        });
        this.tableData = res;
      } catch (error) {
        if (error.data.errInfo) {
          this.$message['warning'](error.data.errInfo)
        } else {
          this.$message['error'](error.message)
        }
      }
    },
    edit(row) {
      this.$refs.form.show(row);
    },
    async handleAble(row) {
      let user = row.user
      let del = row.is_available ? '0':'1'
      try {
        let res = await this.$api.user.handleUserAble({user,del})
        if (res) {
          this.$message["success"]("修改成功");
        }
      } catch (error) {
        if (error.data.errInfo) {
          this.$message['warning'](error.data.errInfo)
        } else {
          this.$message['error'](error.message)
        }
      }
    }
  }
};
</script>

<style lang="less" scoped>
.list {
  height: 700px;
}
.el-tag {
  margin-right: 10px;
}
</style>
