<template>
    <div class="graph-container">
        <div ref="stageContainer" class="konva-stage"></div>
        <div class="graph-legend q-pa-sm" v-if="currentTable">
            <div class="text-caption text-grey-5 uppercase tracking-wider">Relationship Graph: {{ currentTable }}</div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch, onUnmounted } from 'vue';
import Konva from 'konva'
import { io } from '../../wailsjs/go/models';

const props = defineProps<{
    currentTable: string,
    width?: number,
    height?: number,
    relationship: any[],
    lineColor?: string
}>()
  
const stageContainer = ref<HTMLDivElement>()
let stage: Konva.Stage | null = null
let layer: Konva.Layer = new Konva.Layer();
let gridLayer: Konva.Layer = new Konva.Layer()
const containerWidth = ref(800)
const containerHeight = ref(600)


const colors = {
    current: '#3b82f6', // blue-500
    related: '#4b5563', // grey-600
    line: props.lineColor || '#6366f1', // indigo-500
    text: '#ffffff',
    bg: '#111827'
}

const currentTableRelationships = computed(() => {
    return props.relationship.filter(r =>
        r.source.entity === props.currentTable ||
        r.target.entity === props.currentTable
    )
})
const selectedTable = ref<string | null>(props.currentTable);
function drawGridLins() {
    const gridSize = 40;
    const gridColor = 'rgba(255, 255, 255, 0.03)';
    for (let i = 0; i < containerWidth.value / gridSize; i++) {
        gridLayer?.add(new Konva.Line({
            points: [i * gridSize, 0, i * gridSize, containerHeight.value],
            stroke: gridColor,
            strokeWidth: 1,
        }));
    }
    for (let j = 0; j < containerHeight.value / gridSize; j++) {
        gridLayer?.add(new Konva.Line({
            points: [0, j * gridSize, containerWidth.value, j * gridSize],
            stroke: gridColor,
            strokeWidth: 1,
        }));
    }
    gridLayer?.draw();
}

function drawLines() {
    const rels = currentTableRelationships.value;
    const relatedEntities = new Set<string>();
    rels.forEach(r => {
        if (r.source.entity !== props.currentTable) relatedEntities.add(r.source.entity);
        if (r.target.entity !== props.currentTable) relatedEntities.add(r.target.entity);
    });

    const target = layer.findOne<Konva.Group>(`#${selectedTable.value}`);
    const renderedNodes = layer?.find('Group');

    if (!renderedNodes || !target) {
        return;
    }

    for (const element of renderedNodes) {
        const node = element as Konva.Group;
        const nodeName = node.id();

        if (nodeName === selectedTable.value) {
            continue;
        }

        const nodePosition = node.absolutePosition();
        let { x, y } = target.absolutePosition();
        y += radius;

        const startAt = { x, y };
        const endAt = { x: nodePosition.x, y: nodePosition.y - radius };

        // Get all previously rendered nodes (before current one) to avoid
        const currentRenderedPositions = renderedPositions.reduce((acc, v) => {
            if (v.name === nodeName) {
                return acc; // Stop adding when current node is found
            }
            acc.push(v);
            return acc;
        }, [] as typeof renderedPositions);

        // Calculate path avoiding obstacles
        const points = calculateAvoidancePathWithObstacles(
            startAt.x, startAt.y,
            endAt.x, endAt.y,
            currentRenderedPositions,
            radius
        );

        const line = new Konva.Arrow({
            points: points,
            stroke: colors.line,
            strokeWidth: 2,
            tension: 0.3, // Makes the line smooth if multiple points
            lineCap: 'round',
            lineJoin: 'round',
            id: `${selectedTable.value}-TO-${node.id()}`
        });

        layer?.add(line);
    }

    layer?.draw();
}

// Function to calculate path that avoids obstacles
function calculateAvoidancePathWithObstacles(
    startX: number,
    startY: number,
    endX: number,
    endY: number,
    obstacles: typeof renderedPositions,
    radius: number
): number[] {

    // Check if direct line has collisions
    if (!hasCollisions(startX, startY, endX, endY, obstacles, radius)) {
        return [startX, startY, endX, endY];
    }

    // Calculate path with obstacle avoidance
    const pathPoints = [startX, startY];

    // Sort obstacles by distance to line (closest first)
    const sortedObstacles = obstacles.sort((a, b) => {
        const distA = distanceToLine(startX, startY, endX, endY, a.x, a.y);
        const distB = distanceToLine(startX, startY, endX, endY, b.x, b.y);
        return distA - distB;
    });

    let currentX = startX;
    let currentY = startY;

    for (const obstacle of sortedObstacles) {
        if (lineCircleCollision(currentX, currentY, endX, endY, obstacle.x, obstacle.y, radius)) {
            // Need to go around this obstacle
            const avoidancePoint = calculateAvoidancePoint(
                currentX, currentY,
                endX, endY,
                obstacle.x, obstacle.y,
                radius
            );

            pathPoints.push(avoidancePoint.x, avoidancePoint.y);
            currentX = avoidancePoint.x;
            currentY = avoidancePoint.y;
        }
    }

    // Add final destination
    pathPoints.push(endX, endY);

    return pathPoints;
}

// Calculate a point to go around an obstacle
function calculateAvoidancePoint(
    startX: number,
    startY: number,
    endX: number,
    endY: number,
    obstacleX: number,
    obstacleY: number,
    radius: number
): { x: number; y: number } {

    // Vector from start to end
    const lineVecX = endX - startX;
    const lineVecY = endY - startY;

    // Vector from start to obstacle
    const toObstacleX = obstacleX - startX;
    const toObstacleY = obstacleY - startY;

    // Project obstacle onto the line
    const lineLength = Math.sqrt(lineVecX * lineVecX + lineVecY * lineVecY);
    if (lineLength === 0) return { x: startX, y: startY };

    const unitLineX = lineVecX / lineLength;
    const unitLineY = lineVecY / lineLength;

    const projection = toObstacleX * unitLineX + toObstacleY * unitLineY;

    // Point on line closest to obstacle
    const closestOnLineX = startX + unitLineX * projection;
    const closestOnLineY = startY + unitLineY * projection;

    // Direction to move away from obstacle (perpendicular to line)
    const perpX = -unitLineY;
    const perpY = unitLineX;

    // Calculate avoidance direction (away from obstacle)
    const toObstacleDistX = obstacleX - closestOnLineX;
    const toObstacleDistY = obstacleY - closestOnLineY;

    const dotProduct = toObstacleDistX * perpX + toObstacleDistY * perpY;
    const avoidanceDirectionX = dotProduct > 0 ? perpX : -perpX;
    const avoidanceDirectionY = dotProduct > 0 ? perpY : -perpY;

    // Avoidance point at safe distance
    const safeDistance = radius + 30; // Extra buffer

    return {
        x: closestOnLineX + avoidanceDirectionX * safeDistance,
        y: closestOnLineY + avoidanceDirectionY * safeDistance
    };
}

// Check if line collides with any obstacles
function hasCollisions(
    startX: number,
    startY: number,
    endX: number,
    endY: number,
    obstacles: typeof renderedPositions,
    radius: number
): boolean {

    for (const obstacle of obstacles) {
        if (lineCircleCollision(startX, startY, endX, endY, obstacle.x, obstacle.y, radius)) {
            return true;
        }
    }
    return false;
}

// Calculate distance from point to line segment
function distanceToLine(
    startX: number,
    startY: number,
    endX: number,
    endY: number,
    pointX: number,
    pointY: number
): number {

    const A = pointX - startX;
    const B = pointY - startY;
    const C = endX - startX;
    const D = endY - startY;

    const dot = A * C + B * D;
    const lenSq = C * C + D * D;

    if (lenSq === 0) return Math.sqrt(A * A + B * B); // Start and end points are same

    let param = dot / lenSq;

    if (param < 0) {
        // Closest point is start point
        return Math.sqrt((pointX - startX) * (pointX - startX) + (pointY - startY) * (pointY - startY));
    } else if (param > 1) {
        // Closest point is end point
        return Math.sqrt((pointX - endX) * (pointX - endX) + (pointY - endY) * (pointY - endY));
    } else {
        // Closest point is on the line segment
        const xx = startX + param * C;
        const yy = startY + param * D;
        return Math.sqrt((pointX - xx) * (pointX - xx) + (pointY - yy) * (pointY - yy));
    }
}

// Line-circle collision detection
function lineCircleCollision(
    x1: number, y1: number,
    x2: number, y2: number,
    cx: number, cy: number,
    radius: number
): boolean {
    // Vector from start to end
    const dx = x2 - x1;
    const dy = y2 - y1;

    // Vector from start to circle center
    const fx = x1 - cx;
    const fy = y1 - cy;

    // Quadratic coefficients for solving: |P(t) - C|² = r²
    // P(t) = (x1 + t*dx, y1 + t*dy), t ∈ [0, 1]
    const a = dx * dx + dy * dy;
    const b = 2 * (fx * dx + fy * dy);
    const c = fx * fx + fy * fy - radius * radius;

    // Discriminant
    const discriminant = b * b - 4 * a * c;

    // No intersection
    if (discriminant < 0) return false;

    // One or two intersections
    const sqrtDisc = Math.sqrt(discriminant);
    const t1 = (-b - sqrtDisc) / (2 * a);
    const t2 = (-b + sqrtDisc) / (2 * a);

    // Normalize t to [0, 1] — only consider intersections on the line *segment*
    const t1InSegment = t1 >= 0 && t1 <= 1;
    const t2InSegment = t2 >= 0 && t2 <= 1;

    // If both t1 and t2 are in [0,1], the segment passes *through* the circle (chord) → collision
    if (t1InSegment && t2InSegment) {
        // Special case: if t1 === t2 (tangent), it's just touching — NOT a collision
        if (Math.abs(t1 - t2) < 1e-6) {
            return false; // Tangent → allowed
        }
        return true; // Two distinct intersection points → collision
    }

    // If only one intersection (and not tangent), check if it's an endpoint contact
    if (t1InSegment !== t2InSegment) {
        const t = t1InSegment ? t1 : t2;
        // Check if intersection is *exactly* at start or end point (within epsilon)
        const eps = 1e-5;
        const isStart = Math.abs(t) < eps;
        const isEnd = Math.abs(t - 1) < eps;

        if (isStart || isEnd) {
            return false; // Endpoint contact → allowed (this is a valid connection)
        }
        return true; // Intersection in middle of segment → collision
    }

    return false; // No meaningful interior intersection
}

let radius = 45;
let spacing = 50;
let renderedPositions = [] as { x: number, y: number, name: string }[]

function drawEntities(relationship: any[]) {
    let x = 100;
    let y = 100;

    const renderedColumn = [] as string[]

    for (let index = 0; index < relationship.length; index++) {
        const entity = relationship[index];
        const name = entity.name
        if (renderedColumn.includes(name)) {
            continue
        } else {
            renderedColumn.push(name)
        }
        const isCurrent = name === props.currentTable;

        const group = new Konva.Group({
            x: x,
            y: y,
            id: entity.source.entity,
            draggable: true
        })



        const rect = new Konva.Circle({
            radius: radius,
            fill: isCurrent ? colors.current : colors.related,
            cornerRadius: 8,
            shadowBlur: 10,
            shadowOpacity: 0.3,
            stroke: isCurrent ? '#60a5fa' : '#94a3b8',
            strokeWidth: isCurrent ? 2 : 1
        });

        const textLabel = entity.source.entity.length > 20 ? entity.source.entity.slice(0, 20 - 3) + '...' : entity.source.entity
        const text = new Konva.Text({
            text: textLabel,
            fontSize: 10,
            fontStyle: 'bold',
            fill: colors.text,
            width: radius * 2,
            padding: 10,
            align: 'center',
            ellipsis: true,
            offsetX: radius,
            offsetY: radius * 0.35
        });

        group.on('mouseover', () => {
            document.body.style.cursor = 'pointer';
            rect.scale({ x: 1.1, y: 1.1 });
            layer?.batchDraw();
        });

        group.on('mouseout', () => {
            document.body.style.cursor = 'default';
            rect.scale({ x: 1, y: 1 });
            layer?.batchDraw();
        });

        group.add(rect);
        group.add(text);
        layer?.add(group);
        if (x > containerWidth.value - (spacing + radius * 2)) {
            x = 100;
            y += radius * 2 + spacing;
        } else {
            x += radius * 2 + spacing;
        }

        renderedPositions.push({ x, y, name })
    }
}


function drawGraph() {
    if (!layer || !stage || !props.currentTable) return;

    layer.destroyChildren();

    // Draw Grid

    const rels = currentTableRelationships.value;
    const relatedEntities = new Set<string>();
    rels.forEach(r => {
        if (r.source.entity !== props.currentTable) relatedEntities.add(r.source.entity);
        if (r.target.entity !== props.currentTable) relatedEntities.add(r.target.entity);
    });

    const centerX = containerWidth.value / 2;
    const centerY = containerHeight.value / 2;
    const radius = Math.min(centerX, centerY) * 0.7;

    const nodes: Record<string, { x: number, y: number }> = {};
    nodes[props.currentTable] = { x: centerX, y: centerY };

    const relatedArray = Array.from(relatedEntities);
    relatedArray.forEach((entity, i) => {
        const angle = (i / relatedArray.length) * Math.PI * 2;
        nodes[entity] = {
            x: centerX + radius * Math.cos(angle),
            y: centerY + radius * Math.sin(angle)
        };
    });

    // Draw Lines
    rels.forEach(rel => {
        const start = nodes[rel.source.entity];
        const end = nodes[rel.target.entity];

        if (start && end) {
            const arrow = new Konva.Arrow({
                points: [start.x, start.y, end.x, end.y],
                pointerLength: 10,
                pointerWidth: 10,
                fill: colors.line,
                stroke: colors.line,
                strokeWidth: 2,
                opacity: 0.6,
                tension: 0.5
            });

            // Label for the relationship
            const labelText = `${rel.type}\n(${rel.source.field} → ${rel.target.field})`;
            const label = new Konva.Text({
                text: labelText,
                x: (start.x + end.x) / 2,
                y: (start.y + end.y) / 2 - 20,
                fontSize: 9,
                fill: '#94a3b8',
                align: 'center',
                wrap: 'none'
            });
            label.offsetX(label.width() / 2);

            layer?.add(arrow);
            layer?.add(label);
        }
    });

    // Draw Nodes
    Object.entries(nodes).forEach(([name, pos]) => {
        const isCurrent = name === props.currentTable;

        const group = new Konva.Group({
            x: pos.x,
            y: pos.y,
            draggable: true
        });

        const rect = new Konva.Circle({
            width: 120,
            height: 35,
            offsetX: 60,
            offsetY: 20,
            fill: isCurrent ? colors.current : colors.related,
            cornerRadius: 8,
            shadowBlur: 10,
            shadowOpacity: 0.3,
            stroke: isCurrent ? '#60a5fa' : '#94a3b8',
            strokeWidth: isCurrent ? 2 : 1
        });

        const text = new Konva.Text({
            text: name,
            fontSize: 10,
            fontStyle: 'bold',
            fill: colors.text,
            width: 120,
            padding: 10,
            align: 'center',
            offsetX: 60,
            offsetY: 20
        });

        group.on('mouseover', () => {
            document.body.style.cursor = 'pointer';
            rect.scale({ x: 1.1, y: 1.1 });
            layer?.batchDraw();
        });

        group.on('mouseout', () => {
            document.body.style.cursor = 'default';
            rect.scale({ x: 1, y: 1 });
            layer?.batchDraw();
        });

        group.add(rect);
        group.add(text);
        layer?.add(group);
    });

    layer.draw();
}

function updateSize() {
    if (stageContainer.value) {
        containerWidth.value = stageContainer.value.offsetWidth || 800;
        containerHeight.value = 600;
        stage?.width(containerWidth.value);
        stage?.height(containerHeight.value);
        stage?.draw()
    }
}

onMounted(() => {
    if (!stageContainer.value) return;

    containerWidth.value = stageContainer.value.offsetWidth || 800;

    stage = new Konva.Stage({
        container: stageContainer.value,
        width: containerWidth.value,
        height: containerHeight.value
    });


    stage.add(gridLayer)
    drawGridLins()
    stage.add(layer);

    drawEntities(currentTableRelationships.value);

    drawLines()
    //  drawGraph();

    window.addEventListener('resize', updateSize);
});

onUnmounted(() => {
    window.removeEventListener('resize', updateSize);
});

watch(() => [props.currentTable, props.relationship], () => {
    drawGraph();
}, { deep: true });

</script>

<style scoped>
.graph-container {
    width: 100%;
    height: 600px;
    background: #0d0d0d;
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.05);
    position: relative;
    overflow: hidden;
}

.konva-stage {
    width: 100%;
    height: 100%;
}

.graph-legend {
    position: absolute;
    top: 10px;
    left: 10px;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(4px);
    border-radius: 4px;
    pointer-events: none;
}

.tracking-wider {
    letter-spacing: 0.05em;
}
</style>