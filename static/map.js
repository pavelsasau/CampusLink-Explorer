
/*
 * Projet CampusLink Explorer
 *
 * Selon les consignes du projet, JavaScript ne doit pas être
 * utilisé pour le traitement des données.
 *
 * Une petite quantité de JavaScript a été autorisée
 * par notre enseignant pour améliorer l'utilisation de la carte.
 *
 * Ce fichier sert uniquement au zoom et au déplacement.
 * Les données et les API restent gérées par Go.
 */

document.addEventListener("DOMContentLoaded", async () => {

    // Trouver les éléments | Найти элементы
    const viewport = document.getElementById("map-viewport");
    const zoomIn = document.getElementById("zoom-in");
    const zoomOut = document.getElementById("zoom-out");
    const zoomReset = document.getElementById("zoom-reset");
    const zoomLevel = document.getElementById("zoom-level");

    if (!viewport) return;

    // Charger la carte depuis Go | Загрузить карту через Go
    let svg;

    try {
        const response = await fetch("/map.svg");

        if (!response.ok) {
            throw new Error("HTTP " + response.status);
        }

        const svgText = await response.text();

        const parser = new DOMParser();
        const svgDocument = parser.parseFromString(
            svgText,
            "image/svg+xml"
        );

        if (svgDocument.querySelector("parsererror")) {
            throw new Error("SVG invalide");
        }

        svg = document.importNode(svgDocument.documentElement, true);

        viewport.replaceChildren(svg);

    } catch (error) {
        console.error(error);
        viewport.textContent = "Impossible de charger la carte.";
        return;
    }

    // Paramètres du zoom | Параметры масштаба
    let scale = 1;
    const minScale = 1;
    const maxScale = 6;

    // Position de la carte | Положение карты
    let offsetX = 0;
    let offsetY = 0;

    // Paramètres du déplacement | Параметры перемещения
    let isPointerDown = false;
    let isDragging = false;

    let startX = 0;
    let startY = 0;
    let initialX = 0;
    let initialY = 0;

    // Limiter les déplacements | Ограничить перемещение
    function limitPosition() {
        const width = viewport.clientWidth;
        const height = viewport.clientHeight;

        offsetX = Math.min(0, Math.max(width * (1 - scale), offsetX));
        offsetY = Math.min(0, Math.max(height * (1 - scale), offsetY));
    }

    // Mettre à jour la carte | Обновить карту
    function updateMap() {
        limitPosition();

        svg.style.transform =
            `translate(${offsetX}px, ${offsetY}px) scale(${scale})`;

        if (zoomLevel) {
            zoomLevel.textContent = Math.round(scale * 100) + "%";
        }
    }

    // Changer le zoom | Изменить масштаб
    function changeZoom(newScale, mouseX, mouseY) {
        newScale = Math.max(minScale, Math.min(maxScale, newScale));

        const ratio = newScale / scale;

        offsetX = mouseX - (mouseX - offsetX) * ratio;
        offsetY = mouseY - (mouseY - offsetY) * ratio;

        scale = newScale;

        updateMap();
    }

    // Zoom avec la molette | Масштабирование колёсиком
    viewport.addEventListener("wheel", (event) => {
        event.preventDefault();

        const rect = viewport.getBoundingClientRect();

        const mouseX = event.clientX - rect.left;
        const mouseY = event.clientY - rect.top;

        const factor = event.deltaY < 0 ? 1.15 : 1 / 1.15;

        changeZoom(scale * factor, mouseX, mouseY);

    }, { passive: false });

    // Commencer le déplacement | Начать перетаскивание
    viewport.addEventListener("pointerdown", (event) => {
        if (event.button !== 0) return;

        isPointerDown = true;
        isDragging = false;

        startX = event.clientX;
        startY = event.clientY;

        initialX = offsetX;
        initialY = offsetY;
    });

    // Déplacer la carte | Перемещать карту
    viewport.addEventListener("pointermove", (event) => {
        if (!isPointerDown) return;

        const dx = event.clientX - startX;
        const dy = event.clientY - startY;

        if (!isDragging && Math.hypot(dx, dy) > 5) {
            isDragging = true;
            viewport.classList.add("dragging");
            viewport.setPointerCapture(event.pointerId);
        }

        if (!isDragging) return;

        event.preventDefault();

        offsetX = initialX + dx;
        offsetY = initialY + dy;

        updateMap();
    });

    // Arrêter le déplacement | Закончить перетаскивание
    function stopDragging(event) {
        isPointerDown = false;

        viewport.classList.remove("dragging");

        if (viewport.hasPointerCapture(event.pointerId)) {
            viewport.releasePointerCapture(event.pointerId);
        }

        setTimeout(() => {
            isDragging = false;
        }, 0);
    }

    viewport.addEventListener("pointerup", stopDragging);
    viewport.addEventListener("pointercancel", stopDragging);

    // Empêcher le clic après un déplacement | Запретить клик после перетаскивания
    viewport.addEventListener("click", (event) => {
        if (isDragging) {
            event.preventDefault();
            event.stopPropagation();
        }
    }, true);

    // Bouton + | Кнопка увеличения
    if (zoomIn) {
        zoomIn.addEventListener("click", () => {
            changeZoom(scale * 1.5, viewport.clientWidth / 2, viewport.clientHeight / 2);
        });
    }

    // Bouton - | Кнопка уменьшения
    if (zoomOut) {
        zoomOut.addEventListener("click", () => {
            changeZoom(scale / 1.5, viewport.clientWidth / 2, viewport.clientHeight / 2);
        });
    }

    // Réinitialiser | Сбросить масштаб
    if (zoomReset) {
        zoomReset.addEventListener("click", () => {
            scale = 1;
            offsetX = 0;
            offsetY = 0;
            updateMap();
        });
    }

    // Adapter la carte | Адаптировать карту
    window.addEventListener("resize", updateMap);

    updateMap();
});
