<template>
  <div>
    <el-card shadow="never">
      <div slot="header" class="clearfix">
        签名证书列表
      </div>
      <div class="list">
        <el-row style="margin-bottom: 20px">
          <el-col>
            <el-button
              icon="el-icon-plus"
              type="primary"
              size="small"
              @click="createServer"
              >创建签名证书
            </el-button>
          </el-col>
        </el-row>
        <el-table
          :data="tableData"
          height="600"
          border
          stripe
          v-loading="loading"
          element-loading-text="创建中"
        >
          <el-table-column prop="Dn" label="证书名称"></el-table-column>
          <el-table-column prop="Created" label="创建时间"></el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-button
                size="mini"
                type="primary"
                icon="el-icon-eidt"
                @click="download(scope.row)"
                >下载</el-button
              >
            </template>
          </el-table-column>
        </el-table>

        <el-pagination
          background
          :page-sizes="[5, 10, 50, 100]"
          :page-size="query.rows"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="query.page"
        >
        </el-pagination>
      </div>
    </el-card>

    <FormModel ref="form" @handleLoading="handleLoading" />
  </div>
</template>

<script>
import FormModel from "./components/form";
import { zipDownload } from "@/utils/download";

export default {
  name: "Server",
  components: {
    FormModel
  },
  data() {
    return {
      tableData: [],
      loading: false,
      total: 0,
      query: {
        page: 1,
        rows: 10
      }
    };
  },
  created() {
    this.getSign();
  },
  methods: {
    createServer() {
      this.$refs.form.show();
    },
    handleLoading(data) {
      this.loading = data;
      if (!data) {
        this.getSign();
      }
    },
    async getSign() {
      try {
        let res = await this.$api.sign.GetSignList(this.query);
        if (res) {
          this.total = res.total;
          if (res.list) {
            res.list.forEach(item => {
              item.Created = this.$moment(item.Created).format(
                "YYYY-MM-DD HH:mm:ss"
              );
            });
          }
        }
        this.tableData = res.list;
      } catch (error) {
        if (error.data.errInfo) {
          this.$message["warning"](error.data.errInfo);
        } else {
          this.$message["error"](error.message);
        }
      }
    },
    async download(row) {
      let id = row.Id;
      try {
        let res = await this.$api.sign.DownloadCA({ id });
        zipDownload(res);
      } catch (error) {
        if (error.data.errInfo) {
          this.$message["warning"](error.data.errInfo);
        } else {
          this.$message["error"](error.message);
        }
      }
    },
    handleSizeChange(newSize) {
      this.query.rows = newSize;
      this.getSign();
    },
    handleCurrentChange(newPage) {
      this.query.page = newPage;
      this.getSign();
    }
  }
};
</script>

<style lang="less" scoped>
.list {
  height: 74vh;
}
</style>
