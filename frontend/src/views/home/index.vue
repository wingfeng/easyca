<template>
  <div>
    <el-card shadow="never" style="margin-bottom: 20px">
      <Welcome :name="user.userInfo.Name" />
    </el-card>
    <el-card shadow="never">
      <div slot="header" class="clearfix">
        <span>根证书信息</span>
        <el-button
          v-if="haveRoot"
          style="float: right;"
          type="primary"
          size="mini"
          icon="el-icon-download"
          @click="downloadBefore"
          >下载根证书</el-button
        >
      </div>
      <div class="root-ca">
        <div v-if="haveRoot">
          <el-row>
            <el-col>
              <span>版本：</span>
              <span>V{{ rootInfo.Version }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>序列号：</span>
              <span>{{ rootInfo.SerialNumber }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>颁发者：</span>
              <span>{{ rootInfo.Issuer.CommonName }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>使用者：</span>
              <span>{{ rootInfo.Subject.CommonName }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>有效期：</span>
              <span>{{ rootInfo.NotBefore }} ~ {{ rootInfo.NotAfter }}</span>
            </el-col>
          </el-row>
          <!-- <el-row>
            <el-col
              ><el-button
                type="primary"
                size="mini"
                @click="download"
                icon="el-icon-download"
                >下载根证书</el-button
              ></el-col
            >
          </el-row> -->
        </div>
        <div
          v-else
          class="root-ca-create"
          v-loading="loading"
          element-loading-text="创建中"
        >
          <el-link type="info" @click="createRootCA"
            >暂无根证书，点击进行创建</el-link
          >
        </div>
      </div>
    </el-card>
    <FormModel ref="form" @handleLoading="handleLoading" />
  </div>
</template>

<script>
import Welcome from "@/components/welcome";
import FormModel from "./components/form";
import { mapGetters } from "vuex";
import { textDownload } from "@/utils/download";
import { debounce } from "@/utils/common";

export default {
  name: "home",
  components: {
    Welcome,
    FormModel
  },
  data() {
    return {
      loading: false,
      haveRoot: false,
      rootInfo: null
    };
  },
  computed: {
    ...mapGetters(["user"])
  },
  created() {
    this.getRootInfo();
  },
  methods: {
    createRootCA() {
      this.$refs.form.show();
    },

    handleLoading(data) {
      this.loading = data;
      if (!data) {
        this.getRootInfo();
      }
    },

    async getRootInfo() {
      try {
        let res = await this.$api.root.RootInfo();
        if (res) {
          this.haveRoot = true;
          this.rootInfo = res;
          this.rootInfo.NotBefore = this.$moment(
            this.rootInfo.NotBefore
          ).format("YYYY-MM-DD HH:mm:ss");
          this.rootInfo.NotAfter = this.$moment(this.rootInfo.NotAfter).format(
            "YYYY-MM-DD HH:mm:ss"
          );
        }
      } catch (error) {
        this.haveRoot = false;
        if (error.data.errInfo) {
          this.$message['warning'](error.data.errInfo)
        } else {
          this.$message['error'](error.message)
        }
      }
    },

    downloadBefore:debounce(function() {
      this.download()
    }, 500),
    async download() {
      try {
        let res = await this.$api.root.DownloadCA();
        textDownload(res,"根证书");
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
@topOffset: 50%;
@leftOffset: 50%;

.root-ca {
  position: relative;
  height: 60vh;
  line-height: 40px;

  .root-ca-create {
    position: absolute;
    top: @topOffset;
    left: @leftOffset;
    transform: translate(-@leftOffset, -@topOffset);
  }

  div > .el-row > .el-col > span:first-child {
    width: 70px;
    display: inline-block;
    text-align: right;
    margin-right: 10px;
  }
}
</style>
