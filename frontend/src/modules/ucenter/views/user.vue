<template>
	<cl-crud ref="Crud">
		<cl-row>
			<!-- 刷新按钮 -->
			<cl-refresh-btn />
			<!-- 新增按钮 -->
			<cl-add-btn />
			<!-- 删除按钮 -->
			<cl-multi-delete-btn />
			<cl-flex1 />
			<!-- 关键字搜索 -->
			<cl-search-key />
		</cl-row>

		<cl-row>
			<!-- 数据表格 -->
			<cl-table ref="Table" />
		</cl-row>

		<cl-row>
			<cl-flex1 />
			<!-- 分页控件 -->
			<cl-pagination />
		</cl-row>

		<!-- 新增、编辑 -->
		<cl-upsert ref="Upsert" />
	</cl-crud>
</template>

<script lang="ts" name="ucenter-user" setup>
import { useCrud, useTable, useUpsert } from "@cool-vue/crud";
import { useCool } from "/@/cool";
import { v4 as uuidv4 } from "uuid";

const { service } = useCool();
const shortcuts = [
	{
		text: "7天后",
		value: () => {
			const date = new Date();
			date.setTime(date.getTime() + 3600 * 1000 * 24 * 7);
			return date;
		}
	},
	{
		text: "30天后",
		value: () => {
			const date = new Date();
			date.setTime(date.getTime() + 3600 * 1000 * 24 * 30);
			return date;
		}
	},
	{
		text: "90天后",
		value: () => {
			const date = new Date();
			date.setTime(date.getTime() + 3600 * 1000 * 24 * 90);
			return date;
		}
	},
	{
		text: "180天后",
		value: () => {
			const date = new Date();
			date.setTime(date.getTime() + 3600 * 1000 * 24 * 180);
			return date;
		}
	},
	{
		text: "365天后",
		value: () => {
			const date = new Date();
			date.setTime(date.getTime() + 3600 * 1000 * 24 * 365);
			return date;
		}
	}
];

// cl-upsert 配置
const Upsert = useUpsert({
	items: [
		{ label: "名称", prop: "name", required: true, component: { name: "el-input" } },
		{ label: "UserToken", prop: "token", required: true, component: { name: "el-input" } },
		{
			label: "权限",
			prop: "permis",
			required: true,
			component: { name: "el-input", props: { type: "textarea", rows: 7 ,placeholder: "请输入权限,按英文逗号隔开，例如：grok-free,claude-free"} }
		},
		{
			label: "过期时间",
			prop: "expire_time",
			component: {
				name: "el-date-picker",
				props: { type: "datetime", valueFormat: "YYYY-MM-DD HH:mm:ss", shortcuts }
			},
			required: true
		},
		{
			label: "状态",
			prop: "status",
			component: {
				name: "el-switch",
				props: {
					activeValue: 1,
					inactiveValue: 0
				}
			}
		}
	],
	onOpened(data) {
		// 自动生成uuid 作为token
		if (!data.token) {
			data.token = uuidv4();
		}
	}
});

// cl-table 配置
const Table = useTable({
	columns: [
		{ type: "selection" },
		{ label: "id", prop: "id",sortable: true },
		{ label: "创建时间", prop: "createTime",sortable: true },
		{ label: "更新时间", prop: "updateTime",sortable: true },
		{ label: "名称", prop: "name",sortable: true },
		{ label: "Token", prop: "token" ,sortable: true},
		{ label: "权限", prop: "permis",sortable: true },
		{ label: "邮箱", prop: "email",sortable: true },
		{ label: "过期时间", prop: "expire_time",sortable: true },
		{ label: "状态", prop: "status", component: { name: "cl-switch" },sortable: true },
		{ type: "op", buttons: ["edit", "delete"] }
	]
});

// cl-crud 配置
const Crud = useCrud(
	{
		service: service.ucenter.user
	},
	(app) => {
		app.refresh();
	}
);
</script>
