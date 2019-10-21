<template>
    <div class="sync-set-layout" v-bkloading="{ isLoading: $loading('diffTemplateAndInstances') }">
        <template v-if="noInfo">
            <div class="no-content">
                <img src="../../assets/images/no-content.png" alt="no-content">
                <p>{{$t('无集群模板更新信息')}}</p>
                <bk-button theme="primary" @click="handleGoback">{{$t('返回')}}</bk-button>
            </div>
        </template>
        <template v-else-if="isLatestInfo">
            <div class="no-content">
                <img src="../../assets/images/latest-data.png" alt="no-content">
                <p>{{$t('最新集群模板信息')}}</p>
                <bk-button theme="primary" @click="handleGoback">{{$t('返回')}}</bk-button>
            </div>
        </template>
        <template v-else-if="diffList.length">
            <div class="title clearfix">
                <div class="tips fl">
                    <p class="mr10">{{$t('请确认以下模板修改信息')}}：</p>
                    <span class="mr30">
                        <i class="dot"></i>
                        {{$t('新增模块')}}
                    </span>
                    <span>
                        <i class="dot red"></i>
                        {{$t('删除模块')}}
                    </span>
                </div>
                <bk-checkbox class="expand-all fr"
                    v-if="!isSingleSync"
                    v-model="expandAll"
                    @change="handleExpandAll">
                    {{$t('全部展开')}}
                </bk-checkbox>
            </div>
            <div class="instance-list">
                <set-instance class="instance-item"
                    ref="setInstance"
                    v-for="(instance, index) in diffList"
                    :key="instance.bk_set_id"
                    :instance="instance"
                    :icon-close="diffList.length > 1"
                    :expand="index === 0"
                    @close="handleCloseInstance">
                </set-instance>
            </div>
            <div class="footer">
                <span style="display: inlink-block;"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.U_TOPO),
                        auth: [$OPERATION.U_TOPO]
                    }">
                    <bk-button class="mr10"
                        theme="primary"
                        :disabled="!$isAuthorized($OPERATION.U_TOPO)"
                        @click="handleConfirmSync">
                        {{$t('确认同步')}}
                    </bk-button>
                </span>
                <bk-button class="mr10" @click="handleGoback">{{$t('取消')}}</bk-button>
                <span v-if="!isSingleSync">{{$tc('已选集群实例', setInstancesId.length, { count: setInstancesId.length })}}</span>
            </div>
        </template>
    </div>
</template>

<script>
    import { MENU_BUSINESS_SET_TEMPLATE, MENU_BUSINESS_SERVICE_TOPOLOGY } from '@/dictionary/menu-symbol'
    import setInstance from './set-instance'
    export default {
        components: {
            setInstance
        },
        data () {
            return {
                expandAll: false,
                diffList: [],
                noInfo: false,
                isLatestInfo: false,
                templateName: ''
            }
        },
        computed: {
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            setTemplateId () {
                return this.$route.params['setTemplateId']
            },
            setInstancesId () {
                const id = `${this.business}_${this.setTemplateId}`
                let syncIdMap = this.$store.state.setFeatures.syncIdMap
                const sessionSyncIdMap = sessionStorage.getItem('setSyncIdMap')
                if (!Object.keys(syncIdMap).length && sessionSyncIdMap) {
                    syncIdMap = JSON.parse(sessionSyncIdMap)
                    this.$store.commit('setFeatures/resetSyncIdMap', syncIdMap)
                }
                return syncIdMap[id] || []
            },
            isSingleSync () {
                return !(this.setInstancesId.length > 1)
            }
        },
        async created () {
            this.getSetTemplateInfo()
            this.getDiffData()
        },
        methods: {
            setBreadcrumbs () {
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('集群模板'),
                    route: { name: MENU_BUSINESS_SET_TEMPLATE }
                }, {
                    label: this.templateName,
                    route: {
                        name: 'setTemplateConfig',
                        params: {
                            mode: 'view',
                            templateId: this.setTemplateId
                        },
                        query: {
                            tab: 'instance'
                        }
                    }
                }, {
                    label: this.$t('同步集群模板')
                }])
            },
            async getSetTemplateInfo () {
                try {
                    const info = await this.$store.dispatch('setTemplate/getSingleSetTemplateInfo', {
                        bizId: this.business,
                        setTemplateId: this.setTemplateId
                    })
                    this.templateName = info.name
                    this.setBreadcrumbs()
                } catch (e) {
                    console.error(e)
                }
            },
            async getDiffData () {
                try {
                    if (!this.setInstancesId.length) {
                        this.diffList = []
                        this.noInfo = true
                        return
                    }
                    this.diffList = await this.$store.dispatch('setSync/diffTemplateAndInstances', {
                        bizId: this.business,
                        setTemplateId: this.setTemplateId,
                        params: {
                            bk_set_ids: this.setInstancesId
                        },
                        config: {
                            requestId: 'diffTemplateAndInstances'
                        }
                    })
                    const changeList = this.diffList.filter(set => {
                        const moduleDiffs = set.module_diffs
                        return moduleDiffs && moduleDiffs.filter(module => module.diff_type !== 'unchanged').length
                    })
                    this.isLatestInfo = !changeList.length
                    this.noInfo = false
                } catch (e) {
                    console.error(e)
                    this.noInfo = true
                }
            },
            async handleConfirmSync () {
                try {
                    await this.$store.dispatch('setSync/syncTemplateToInstances', {
                        bizId: this.business,
                        setTemplateId: this.setTemplateId,
                        params: {
                            bk_set_ids: this.setInstancesId
                        },
                        config: {
                            requestId: 'syncTemplateToInstances'
                        }
                    })
                    this.$success(this.$t('提交同步成功'))
                    this.$router.replace({
                        name: 'setTemplateConfig',
                        params: {
                            templateId: this.setTemplateId,
                            mode: 'view'
                        },
                        query: {
                            tab: 'instance'
                        }
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            handleExpandAll (expand) {
                this.$refs.setInstance.forEach(instance => {
                    instance.localExpand = expand
                })
            },
            handleCloseInstance (id) {
                this.$store.commit('setFeatures/deleteInstancesId', {
                    id: `${this.business}_${this.setTemplateId}`,
                    deleteId: id
                })
                this.diffList = this.diffList.filter(instance => instance.bk_set_id !== id)
            },
            handleGoback () {
                const moduleId = this.$route.params['moduleId']
                if (moduleId) {
                    this.$router.replace({
                        name: MENU_BUSINESS_SERVICE_TOPOLOGY,
                        query: {
                            module: moduleId
                        }
                    })
                } else {
                    this.$router.replace({
                        name: 'setTemplateConfig',
                        params: {
                            templateId: this.setTemplateId,
                            mode: 'view'
                        },
                        query: {
                            tab: 'instance'
                        }
                    })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .sync-set-layout {
        padding: 0 20px;
    }
    .no-content {
        position: absolute;
        top: 50%;
        left: 50%;
        font-size: 16px;
        color: #63656e;
        text-align: center;
        transform: translate(-50%, -50%);
        img {
            width: 130px;
        }
        p {
            padding: 20px 0 30px;
        }
    }
    .tips {
        display: flex;
        align-items: center;
        font-size: 14px;
        color: #63656E;
        .dot {
            display: inline-block;
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background-color: #2DCB56;
            margin-right: 2px;
            &.red {
                background-color: #FF5656;
            }
        }
    }
    .expand-all {
        color: #888991;
    }
    .instance-list {
        padding: 20px 0 0;
        .instance-item {
            margin-bottom: 10px;
        }
    }
    .footer {
        position: sticky;
        bottom: 0;
        display: flex;
        align-items: center;
        padding: 10px 0 20px;
        background-color: #fafbfd;
        > span {
            color: #979BA5;
            font-size: 14px;
        }
    }
</style>

<style lang="scss">
    .set-confirm-sync {
        .bk-dialog-content {
            width: 440px !important;
        }
        .bk-dialog-type-body {
            padding: 2px 24px 5px !important;
        }
        .bk-dialog-type-sub-header {
            padding: 3px 40px 24px !important;
            .header {
                white-space: unset !important;
                text-align: left !important;
            }
        }
        .bk-dialog-footer {
            padding-bottom: 32px !important;
        }
    }
</style>