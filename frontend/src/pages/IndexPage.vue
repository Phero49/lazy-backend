<template>
    <q-page class="q-pa-xl">
        <div class="max-width-1200 q-mx-auto">
            <!-- Welcome Header -->
            <div class="q-mb-xl">
                <h1 class="text-h3 text-white text-weight-bolder q-mb-xs tracking-tight">
                    System <span class="text-primary">Overview</span>
                </h1>
                <p class="text-grey-6 text-subtitle1">
                    Monitor your services and manage backend resources efficiently.
                </p>
            </div>

            <!-- Live Statistics -->
            <div class="row q-col-gutter-lg q-mb-xl">
                <div v-for="(stat, i) in stats" :key="i" class="col-12 col-md-3">
                    <q-card flat class="stat-card bg-dark-soft border-lighter q-pa-lg">
                        <div class="row items-center justify-between q-mb-md">
                            <q-icon :name="stat.icon" :color="stat.color" size="24px"
                                class="bg-primary-opacity q-pa-sm border-radius-sm" />
                            <q-badge :color="stat.badgeColor" rounded class="q-px-sm">{{ stat.status }}</q-badge>
                        </div>
                        <div class="text-h4 text-white text-weight-bold">{{ stat.value }}</div>
                        <div class="text-caption text-grey-7 text-uppercase tracking-wider text-weight-bold">{{
                            stat.label }}</div>
                    </q-card>
                </div>
            </div>

            <!-- Action Grid -->
            <div class="row q-col-gutter-xl">
                <!-- Recent Services -->
                <div class="col-12 col-md-8">
                    <div class="row items-center justify-between q-mb-md">
                        <span class="text-subtitle1 text-white text-weight-bold">Active Service Instances</span>
                        <q-btn flat dense no-caps color="primary" label="View all services"
                            icon-right="bi-chevron-right" />
                    </div>

                    <div class="column q-gutter-y-md">
                        <q-card v-for="i in 3" :key="i" flat
                            class="bg-dark-soft border-lighter card-hover q-pa-md cursor-pointer">
                            <div class="row items-center no-wrap">
                                <div class="q-pa-md bg-dark border-lighter border-radius-md q-mr-lg">
                                    <q-icon name="bi-server" color="primary" size="20px" />
                                </div>
                                <div class="col">
                                    <div class="text-weight-bold text-grey-2">production-instance-0{{ i }}</div>
                                    <div class="text-caption text-grey-7">Endpoint: api.lazybackend.io/v{{ i }} • Last
                                        sync: 4m ago</div>
                                </div>
                                <div class="column items-end">
                                    <div class="text-primary text-weight-bold">99.9% Uptime</div>
                                    <div class="size-10 text-grey-9 text-uppercase">v1.2.{{ i }}</div>
                                </div>
                            </div>
                        </q-card>
                    </div>
                </div>

                <!-- System Health Widget -->
                <div class="col-12 col-md-4">
                    <q-card flat
                        class="bg-primary text-white q-pa-xl border-radius-lg relative-position overflow-hidden">
                        <q-icon name="bi-cpu" class="absolute-bottom-right icon-watermark" />
                        <div class="text-h5 text-weight-bolder q-mb-md">Core Performance</div>
                        <p class="text-body2 text-white-70 q-mb-xl">
                            Go-routines are currently optimized for high-throughput packet processing.
                        </p>
                        <div class="row items-end justify-between">
                            <div class="column">
                                <span class="text-h4 text-weight-bold">2.4ms</span>
                                <span class="text-caption opacity-70">Avg Latency</span>
                            </div>
                            <q-btn unelevated color="white" text-color="primary" label="Optimize"
                                class="text-weight-bold border-radius-md" no-caps />
                        </div>
                    </q-card>
                </div>
            </div>
        </div>
    </q-page>
</template>

<script setup lang="ts">
const stats = [
    { label: 'Total Requests', value: '45.2k', icon: 'bi-activity', color: 'primary', badgeColor: 'green-10', status: '+12%' },
    { label: 'Active Tasks', value: '84', icon: 'bi-check-circle', color: 'primary', badgeColor: 'green-10', status: 'Stable' },
    { label: 'System Errors', value: '0', icon: 'bi-bug-fill', color: 'red-5', badgeColor: 'red-10', status: 'None' },
    { label: 'Connected Nodes', value: '12', icon: 'bi-share-fill', color: 'primary', badgeColor: 'green-10', status: 'Live' },
];
</script>

<style lang="scss" scoped>
.max-width-1200 {
    max-width: 1200px;
}

.tracking-tight {
    letter-spacing: -0.04em;
}

.tracking-wider {
    letter-spacing: 0.1em;
}

.size-10 {
    font-size: 10px;
}

.text-white-70 {
    color: rgba(255, 255, 255, 0.7);
}

.bg-dark {
    background: #050505;
}

.bg-dark-soft {
    background: #0d0d0d;
}

.border-lighter {
    border: 1px solid rgba(255, 255, 255, 0.05);
}

.border-radius-sm {
    border-radius: 8px;
}

.border-radius-md {
    border-radius: 12px;
}

.border-radius-lg {
    border-radius: 20px;
}

.bg-primary-opacity {
    background: rgba(var(--q-primary), 0.1);
}

.stat-card {
    transition: transform 0.2s, border-color 0.2s;

    &:hover {
        transform: translateY(-4px);
        border-color: rgba(var(--q-primary), 0.3);
    }
}



.icon-watermark {
    font-size: 120px;
    right: -20px;
    bottom: -20px;
    opacity: 0.1;
    color: #fff;
}

.opacity-70 {
    opacity: 0.7;
}
</style>
