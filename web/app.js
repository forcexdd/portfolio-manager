const content = document.getElementById('view-port');
const viewTitle = document.getElementById('current-view-title');
const globalPortSelector = document.getElementById('global-portfolio-select');

let appState = {
    currentView: 'analysis',
    selectedPortfolioId: null,
    selectedIndexId: null,
    portfolios: [],
    indexes: []
};

async function api(url, method = 'GET', body = null) {
    const opts = { method, headers: {} };
    if (body) {
        opts.headers['Content-Type'] = 'application/json';
        opts.body = JSON.stringify(body);
    }
    try {
        const res = await fetch(url, opts);
        if (!res.ok) throw new Error(await res.text());
        if (res.status === 204) return null;
        return await res.json();
    } catch (err) {
        alert(`API Error: ${err.message}`);
        throw err;
    }
}

window.addEventListener('hashchange', handleHashRouting);

async function handleHashRouting() {
    const hash = window.location.hash.replace('#', '') || 'analysis';
    if (appState.currentView !== hash) {
        await switchView(hash, false);
    }
}

async function initializeApp() {
    await reloadConfigData();
    
    const savedPortId = localStorage.getItem('selected_portfolio_id');
    if (savedPortId && appState.portfolios.some(p => p.id === parseInt(savedPortId))) {
        appState.selectedPortfolioId = parseInt(savedPortId);
    } else if (appState.portfolios.length > 0) {
        appState.selectedPortfolioId = appState.portfolios[0].id;
    }

    const savedIdxId = localStorage.getItem('selected_index_id');
    if (savedIdxId && appState.indexes.some(i => i.id === parseInt(savedIdxId))) {
        appState.selectedIndexId = parseInt(savedIdxId);
    } else if (appState.indexes.length > 0) {
        appState.selectedIndexId = appState.indexes[0].id;
    }

    renderGlobalSelector();
    
    const initView = window.location.hash.replace('#', '') || 'analysis';
    await switchView(initView, false);
}

async function reloadConfigData() {
    appState.portfolios = await api('/api/portfolios') || [];
    appState.indexes = await api('/api/indexes') || [];
    
    appState.portfolios.sort((a, b) => a.name.localeCompare(b.name));
    appState.indexes.sort((a, b) => a.name.localeCompare(b.name));
}

function renderGlobalSelector() {
    if (appState.portfolios.length === 0) {
        globalPortSelector.innerHTML = `<option value="">No portfolios defined</option>`;
        return;
    }
    globalPortSelector.innerHTML = appState.portfolios.map(p => 
        `<option value="${p.id}" ${p.id === appState.selectedPortfolioId ? 'selected' : ''}>${p.name}</option>`
    ).join('');
}

async function onGlobalPortfolioChange() {
    appState.selectedPortfolioId = parseInt(globalPortSelector.value);
    localStorage.setItem('selected_portfolio_id', appState.selectedPortfolioId);
    if (appState.currentView === 'analysis') {
        await renderAnalysisView();
    } else if (appState.currentView === 'portfolios') {
        await renderPortfoliosView();
    }
}

async function switchView(viewName, updateHash = true) {
    appState.currentView = viewName;
    if (updateHash) window.location.hash = viewName;

    document.querySelectorAll('.nav-item').forEach(el => el.classList.remove('active'));
    const activeBtn = document.getElementById(`btn-${viewName}`);
    if (activeBtn) activeBtn.classList.add('active');

    if (viewName === 'analysis') {
        viewTitle.innerText = "Analysis Panel";
        await renderAnalysisView();
    } else if (viewName === 'portfolios') {
        viewTitle.innerText = "Portfolios Management";
        await renderPortfoliosView();
    } else if (viewName === 'assets') {
        viewTitle.innerText = "System Assets Directory";
        await renderAssetsView();
    }
}

async function renderAnalysisView() {
    if (!appState.selectedPortfolioId) {
        content.innerHTML = `
            <div class="card" style="text-align: center; padding: 48px;">
                <h3>No active portfolio selected</h3>
                <p class="text-muted" style="margin-top: 8px;">Please define or select a portfolio to analyze.</p>
                <button class="primary" style="margin-top: 16px;" onclick="switchView('portfolios')">Go to Portfolios</button>
            </div>`;
        return;
    }

    if (appState.indexes.length === 0) {
        content.innerHTML = `
            <div class="card" style="text-align: center; padding: 48px;">
                <h3>No indices available</h3>
                <p class="text-muted">The system is fetching indexing details. Please wait.</p>
            </div>`;
        return;
    }

    if (!appState.selectedIndexId && appState.indexes.length > 0) {
        appState.selectedIndexId = appState.indexes[0].id;
    }

    const analysisData = await api(`/api/analysis?portfolio_id=${appState.selectedPortfolioId}&index_id=${appState.selectedIndexId}`);
    const assets = analysisData.assets || [];
    assets.sort((a, b) => a.name.localeCompare(b.name));

    let absoluteDeviation = 0;
    assets.forEach(a => {
        absoluteDeviation += Math.abs(a.fraction_diff);
    });
    const fidelityScore = Math.max(0, 100 - absoluteDeviation);

    content.innerHTML = `
        <div class="stats-grid">
            <div class="stat-card">
                <span class="stat-label">Total Valuation</span>
                <span class="stat-value">₽ ${analysisData.total_value.toLocaleString('en-US', {minimumFractionDigits: 2, maximumFractionDigits: 2})}</span>
            </div>
            <div class="stat-card">
                <span class="stat-label">Reference Index</span>
                <select id="analysis-index-select" onchange="onAnalysisIndexChange()" style="margin-top: 4px; font-weight: 600;">
                    ${appState.indexes.map(i => `<option value="${i.id}" ${i.id === appState.selectedIndexId ? 'selected' : ''}>${i.name}</option>`).join('')}
                </select>
            </div>
            <div class="stat-card">
                <span class="stat-label">Fidelity Match</span>
                <span class="stat-value" style="color: ${fidelityScore > 80 ? 'var(--emerald)' : 'var(--amber)'}">${fidelityScore.toFixed(1)}%</span>
            </div>
        </div>

        <div class="grid-two-col">
            <div class="data-table-container">
                <table>
                    <thead>
                        <tr>
                            <th>Asset</th>
                            <th>Price</th>
                            <th>Your %</th>
                            <th>Index %</th>
                            <th>Diff %</th>
                            <th>Lot Size</th>
                            <th>Current Qty</th>
                            <th>Target Qty</th>
                            <th>Action</th>
                            <th style="text-align: right; width: 340px;">Execute Trade</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${assets.map(a => {
                            const lotsDelta = Math.round(a.difference / a.lot_size);
                            const diffColor = a.fraction_diff > 0 ? 'var(--emerald)' : a.fraction_diff < 0 ? 'var(--rose)' : 'inherit';
                            
                            let tradeButtonHtml = '—';
                            if (lotsDelta > 0) {
                                tradeButtonHtml = `<button class="primary" style="padding: 4px 8px; font-size: 11px; background-color: var(--emerald); border-color: var(--emerald);" onclick="executeLotTrade(${appState.selectedPortfolioId}, ${a.asset_id}, ${a.current_quantity}, ${lotsDelta * a.lot_size})">Buy ${lotsDelta} Lot(s)</button>`;
                            } else if (lotsDelta < 0) {
                                tradeButtonHtml = `<button class="danger" style="padding: 4px 8px; font-size: 11px;" onclick="executeLotTrade(${appState.selectedPortfolioId}, ${a.asset_id}, ${a.current_quantity}, ${lotsDelta * a.lot_size})">Sell ${Math.abs(lotsDelta)} Lot(s)</button>`;
                            } else if (a.target_quantity === 0 && a.current_quantity > 0) {
                                tradeButtonHtml = `<button class="danger" style="padding: 4px 8px; font-size: 11px;" onclick="executeLotTrade(${appState.selectedPortfolioId}, ${a.asset_id}, ${a.current_quantity}, ${-a.current_quantity})">Liquidate All</button>`;
                            }

                            return `
                                <tr>
                                    <td><strong>${a.name}</strong></td>
                                    <td>₽${a.price.toFixed(2)}</td>
                                    <td>${a.current_fraction.toFixed(2)}%</td>
                                    <td>${a.index_fraction.toFixed(2)}%</td>
                                    <td style="color: ${diffColor}">
                                        ${a.fraction_diff > 0 ? '+' : ''}${a.fraction_diff.toFixed(2)}%
                                    </td>
                                    <td>
                                        <input type="number" class="inline-edit" value="${a.lot_size}" onchange="updateLotSize(${appState.selectedPortfolioId}, ${a.asset_id}, this.value)">
                                    </td>
                                    <td>
                                        <input type="number" class="inline-edit" value="${a.current_quantity}" onchange="updateQuantity(${appState.selectedPortfolioId}, ${a.asset_id}, this.value)">
                                    </td>
                                    <td>${a.target_quantity}</td>
                                    <td><span class="badge badge-${a.action.toLowerCase()}">${a.action}</span></td>
                                    <td>
                                        <div style="display: flex; align-items: center; gap: 8px; justify-content: flex-end;">
                                            ${tradeButtonHtml}
                                            <span style="border-left: 1px solid var(--border); height: 18px;"></span>
                                            <div style="display: flex; align-items: center; gap: 4px;">
                                                <input type="number" id="manual-lots-${a.asset_id}" min="1" value="1" class="inline-edit" style="width: 44px; text-align: center;">
                                                <button class="primary" style="padding: 4px 6px; font-size: 11px; background-color: var(--emerald); border-color: var(--emerald);" 
                                                        onclick="executeManualLots(${appState.selectedPortfolioId}, ${a.asset_id}, ${a.current_quantity}, ${a.lot_size}, 'BUY')">+</button>
                                                <button class="danger" style="padding: 4px 6px; font-size: 11px;" 
                                                        onclick="executeManualLots(${appState.selectedPortfolioId}, ${a.asset_id}, ${a.current_quantity}, ${a.lot_size}, 'SELL')">-</button>
                                            </div>
                                        </div>
                                    </td>
                                </tr>
                            `;
                        }).join('')}
                    </tbody>
                </table>
            </div>

            <div style="display: flex; flex-direction: column; gap: 24px;">
                <div class="card" style="align-self: start; width: 100%;">
                    <h3 style="margin-bottom: 16px;">Quick Add Asset</h3>
                    <input type="text" id="asset-search" class="search-bar" placeholder="Search by asset ticker..." oninput="filterAssetsDropdown()">
                    <select id="quick-asset-select" style="width: 100%; margin-bottom: 16px;"></select>
                    <input id="quick-qty" type="number" value="10" style="width: 100%; margin-bottom: 16px;" placeholder="Quantity">
                    <button class="primary" style="width: 100%;" onclick="submitQuickAdd()">Add to Portfolio</button>
                </div>

                <div class="card" style="align-self: start; width: 100%;">
                    <h3 style="margin-bottom: 16px;">Target Adjustments</h3>
                    ${pendingInstructionsHtml(pendingAdjustments(assets))}
                </div>
            </div>
        </div>`;

    await loadQuickAssetDropdown();
}

async function executeLotTrade(pid, aid, currentQty, difference) {
    const newQty = currentQty + difference;
    if (newQty <= 0) {
        await api(`/api/portfolios/${pid}/assets/${aid}`, 'DELETE');
    } else {
        await api(`/api/portfolios/${pid}/assets`, 'POST', { asset_id: aid, quantity: newQty });
    }
    await renderAnalysisView();
}

async function executeManualLots(pid, aid, currentQty, lotSize, direction) {
    const input = document.getElementById(`manual-lots-${aid}`);
    if (!input) return;
    
    const lots = parseInt(input.value);
    if (isNaN(lots) || lots <= 0) return;
    
    const deltaShares = lots * lotSize;
    let newQty = currentQty;
    if (direction === 'BUY') {
        newQty += deltaShares;
    } else {
        newQty -= deltaShares;
    }
    
    if (newQty <= 0) {
        await api(`/api/portfolios/${pid}/assets/${aid}`, 'DELETE');
    } else {
        await api(`/api/portfolios/${pid}/assets`, 'POST', { asset_id: aid, quantity: newQty });
    }
    await renderAnalysisView();
}

function pendingAdjustments(assets) {
    const pending = assets.filter(a => {
        const lotDiff = Math.round(a.difference / a.lot_size);
        return lotDiff !== 0 || (a.target_quantity === 0 && a.current_quantity > 0);
    });
    pending.sort((a, b) => Math.abs(b.fraction_diff) - Math.abs(a.fraction_diff));
    return pending;
}

function pendingInstructionsHtml(pending) {
    if (pending.length === 0) {
        return `<p class="text-muted">Portfolio matches index precisely. No transactions needed.</p>`;
    }
    let html = `<ul style="list-style: none; padding: 0; display: flex; flex-direction: column; gap: 12px;">`;
    pending.forEach(a => {
        const lotDiff = Math.round(a.difference / a.lot_size);
        let actionLabel = "";
        let tradeDiff = 0;
        let actionColor = "";
        let execBtnClass = "";
        let targetChange = 0;

        if (lotDiff > 0) {
            actionLabel = "BUY";
            tradeDiff = lotDiff;
            actionColor = "var(--emerald)";
            execBtnClass = "primary";
            targetChange = lotDiff * a.lot_size;
        } else if (lotDiff < 0) {
            actionLabel = "SELL";
            tradeDiff = Math.abs(lotDiff);
            actionColor = "var(--rose)";
            execBtnClass = "danger";
            targetChange = lotDiff * a.lot_size;
        } else if (a.target_quantity === 0 && a.current_quantity > 0) {
            actionLabel = "LIQUIDATE";
            tradeDiff = a.current_quantity;
            actionColor = "var(--rose)";
            execBtnClass = "danger";
            targetChange = -a.current_quantity;
        }

        html += `
            <li style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border); padding-bottom: 10px;">
                <div>
                    <span style="font-weight: 600; color: #fff;">${a.name}</span> 
                    <span style="font-size: 11px; color: var(--text-muted);">(impact: ${Math.abs(a.fraction_diff).toFixed(2)}%)</span>
                    <br>
                    <span class="text-muted" style="font-size: 12px;">Lot Size: ${a.lot_size}</span>
                </div>
                <div style="text-align: right; display: flex; align-items: center; gap: 10px;">
                    <div>
                        <span style="color: ${actionColor}; font-weight: 700;">
                            ${actionLabel} ${actionLabel === "LIQUIDATE" ? "" : tradeDiff + " Lot(s)"}
                        </span>
                        <br>
                        <span class="text-muted" style="font-size: 11px;">(${Math.abs(targetChange)} shares)</span>
                    </div>
                    <button class="${execBtnClass}" style="padding: 4px 8px; font-size: 12px; ${actionLabel === "BUY" ? 'background-color: var(--emerald); border-color: var(--emerald);' : ''}" 
                            onclick="executeLotTrade(${appState.selectedPortfolioId}, ${a.asset_id}, ${a.current_quantity}, ${targetChange})">✔</button>
                </div>
            </li>`;
    });
    return html + `</ul>`;
}

async function onAnalysisIndexChange() {
    appState.selectedIndexId = parseInt(document.getElementById('analysis-index-select').value);
    localStorage.setItem('selected_index_id', appState.selectedIndexId);
    await renderAnalysisView();
}

async function loadQuickAssetDropdown() {
    loadedAssets = await api('/api/assets') || [];
    loadedAssets.sort((a, b) => a.name.localeCompare(b.name));
    filterAssetsDropdown();
}

function filterAssetsDropdown() {
    const searchVal = (document.getElementById('asset-search')?.value || '').toUpperCase();
    const dropdown = document.getElementById('quick-asset-select');
    if (!dropdown) return;

    const filtered = loadedAssets.filter(a => a.name.toUpperCase().includes(searchVal));
    dropdown.innerHTML = filtered.map(a => 
        `<option value="${a.id}">${a.name} (₽${a.price.toFixed(2)})</option>`
    ).join('');
}

async function submitQuickAdd() {
    const select = document.getElementById('quick-asset-select');
    const qtyInput = document.getElementById('quick-qty');
    if (!select || !qtyInput) return;

    const assetId = parseInt(select.value);
    const quantity = parseInt(qtyInput.value);

    if (!isNaN(assetId) && !isNaN(quantity) && quantity > 0) {
        await api(`/api/portfolios/${appState.selectedPortfolioId}/assets`, 'POST', { asset_id: assetId, quantity });
        await renderAnalysisView();
    }
}

async function updateQuantity(pid, aid, newQty) {
    const qty = parseInt(newQty);
    if (qty <= 0) {
        await api(`/api/portfolios/${pid}/assets/${aid}`, 'DELETE');
    } else {
        await api(`/api/portfolios/${pid}/assets`, 'POST', { asset_id: aid, quantity: qty });
    }
    await renderAnalysisView();
}

async function updateLotSize(pid, aid, newLot) {
    const lot = parseInt(newLot);
    if (!isNaN(lot) && lot >= 1) {
        await api(`/api/portfolios/${pid}/assets/${aid}/lot`, 'PUT', { lot_size: lot });
        await renderAnalysisView();
    }
}

async function renderPortfoliosView() {
    await reloadConfigData();
    renderGlobalSelector();

    let portfolioAssets = [];
    if (appState.selectedPortfolioId) {
        portfolioAssets = await api(`/api/portfolios/${appState.selectedPortfolioId}/assets`) || [];
        portfolioAssets.sort((a, b) => a.asset.name.localeCompare(b.asset.name));
    }

    content.innerHTML = `
        <div class="grid-two-col">
            <div>
                <div class="card" style="margin-bottom: 24px;">
                    <h3 style="margin-bottom: 16px;">Create New Portfolio</h3>
                    <div class="flex-container">
                        <input id="new-portfolio-name" style="flex-grow: 1;" placeholder="Asset Bundle Name">
                        <button class="primary" onclick="createPortfolio()">Register Bundle</button>
                    </div>
                    <div style="margin-top: 16px; padding-top: 16px; border-top: 1px solid var(--border);">
                        <h4 style="margin-bottom: 8px;">Import Profile</h4>
                        <input type="file" id="import-profile-file" accept=".json" style="width: 100%; max-width: 300px; margin-bottom: 8px;">
                        <br>
                        <button onclick="importPortfolioProfile()">Import JSON Profile</button>
                    </div>
                </div>

                <div class="data-table-container">
                    <div style="padding: 20px 24px; font-weight: 600; font-size: 16px; border-bottom: 1px solid var(--border);">
                        Registered Portfolios
                    </div>
                    <table>
                        <thead>
                            <tr>
                                <th>Portfolio Name</th>
                                <th style="text-align: right;">Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            ${appState.portfolios.map(p => `
                                <tr class="${p.id === appState.selectedPortfolioId ? 'bg-surface-hover' : ''}">
                                    <td><strong>${p.name}</strong></td>
                                    <td style="text-align: right;">
                                        <button onclick="selectAndExportPortfolio(${p.id})">Export</button>
                                        <button onclick="selectAndManagePortfolio(${p.id})">Manage Assets</button>
                                        <button class="danger" onclick="deletePortfolio(${p.id})">Delete</button>
                                    </td>
                                </tr>
                            `).join('')}
                        </tbody>
                    </table>
                </div>
            </div>

            <div>
                <div class="card" style="margin-bottom: 24px;">
                    <h3>Asset Composition</h3>
                    <p class="text-muted" style="margin-bottom: 16px;">Inline inventory override</p>
                    ${portfolioAssets.length === 0 ? '<p class="text-muted">No assets found in current portfolio.</p>' : `
                        <table style="width: 100%;">
                            <thead>
                                <tr>
                                    <th>Asset</th>
                                    <th>Lot Size</th>
                                    <th>Quantity</th>
                                    <th>Action</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${portfolioAssets.map(pa => `
                                    <tr>
                                        <td><strong>${pa.asset.name}</strong></td>
                                        <td>
                                            <input type="number" class="inline-edit" min="1" value="${pa.lot_size}" 
                                                   onchange="updatePortfolioViewLotSize(${appState.selectedPortfolioId}, ${pa.asset.id}, this.value)">
                                        </td>
                                        <td>
                                            <input type="number" class="inline-edit" value="${pa.quantity}" 
                                                   onchange="updatePortfolioViewQty(${appState.selectedPortfolioId}, ${pa.asset.id}, this.value)">
                                        </td>
                                        <td>
                                            <button class="danger" style="padding: 4px 8px;" onclick="deletePortfolioAsset(${appState.selectedPortfolioId}, ${pa.asset.id})">×</button>
                                        </td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    `}
                </div>

                ${appState.selectedPortfolioId ? `
                    <div class="card">
                        <h3>Add Asset to Selection</h3>
                        <input type="text" id="portfolio-asset-search" class="search-bar" placeholder="Search system asset catalog..." oninput="filterPortfolioAssetsDropdown()">
                        <select id="portfolio-quick-asset-select" style="width: 100%; margin-bottom: 16px;"></select>
                        <input id="portfolio-quick-qty" type="number" value="10" style="width: 100%; margin-bottom: 16px;">
                        <button class="primary" style="width: 100%;" onclick="submitPortfolioQuickAdd(${appState.selectedPortfolioId})">Add Asset</button>
                    </div>
                ` : ''}
            </div>
        </div>`;

    if (appState.selectedPortfolioId) {
        await loadPortfolioQuickAssetDropdown();
    }
}

async function loadPortfolioQuickAssetDropdown() {
    loadedAssets = await api('/api/assets') || [];
    filterPortfolioAssetsDropdown();
}

function filterPortfolioAssetsDropdown() {
    const searchVal = (document.getElementById('portfolio-asset-search')?.value || '').toUpperCase();
    const dropdown = document.getElementById('portfolio-quick-asset-select');
    if (!dropdown) return;

    const filtered = loadedAssets.filter(a => a.name.toUpperCase().includes(searchVal));
    dropdown.innerHTML = filtered.map(a => 
        `<option value="${a.id}">${a.name} (₽${a.price.toFixed(2)})</option>`
    ).join('');
}

async function submitPortfolioQuickAdd(pid) {
    const select = document.getElementById('portfolio-quick-asset-select');
    const qtyInput = document.getElementById('portfolio-quick-qty');
    if (!select || !qtyInput) return;

    const assetId = parseInt(select.value);
    const quantity = parseInt(qtyInput.value);

    if (!isNaN(assetId) && !isNaN(quantity) && quantity > 0) {
        await api(`/api/portfolios/${pid}/assets`, 'POST', { asset_id: assetId, quantity: quantity });
        await renderPortfoliosView();
    }
}

async function selectAndManagePortfolio(id) {
    appState.selectedPortfolioId = id;
    localStorage.setItem('selected_portfolio_id', id);
    renderGlobalSelector();
    await renderPortfoliosView();
}

async function selectAndExportPortfolio(id) {
    const profile = await api(`/api/portfolios/${id}/export`);
    const blob = new Blob([JSON.stringify(profile, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${profile.name.toLowerCase().replace(/\s+/g, '_')}_profile.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

async function importPortfolioProfile() {
    const fileInput = document.getElementById('import-profile-file');
    if (!fileInput || !fileInput.files || fileInput.files.length === 0) {
        alert("Please select a valid JSON profile file first.");
        return;
    }
    const file = fileInput.files[0];
    const reader = new FileReader();
    reader.onload = async (e) => {
        try {
            const profile = JSON.parse(e.target.result);
            if (!profile.name || !profile.assets) {
                throw new Error("Invalid profile layout format");
            }
            const p = await api('/api/portfolios/import', 'POST', profile);
            appState.selectedPortfolioId = p.id;
            localStorage.setItem('selected_portfolio_id', p.id);
            await renderPortfoliosView();
        } catch (err) {
            alert(`File Import failed: ${err.message}`);
        }
    };
    reader.readAsText(file);
}

async function createPortfolio() {
    const name = document.getElementById('new-portfolio-name').value;
    if (name) {
        const p = await api('/api/portfolios', 'POST', { name });
        appState.selectedPortfolioId = p.id;
        localStorage.setItem('selected_portfolio_id', p.id);
        await renderPortfoliosView();
    }
}

async function deletePortfolio(id) {
    if (confirm('Are you sure you want to permanently delete this portfolio?')) {
        await api(`/api/portfolios/${id}`, 'DELETE');
        if (appState.selectedPortfolioId === id) {
            appState.selectedPortfolioId = null;
            localStorage.removeItem('selected_portfolio_id');
        }
        await renderPortfoliosView();
    }
}

async function updatePortfolioViewQty(pid, aid, newQty) {
    const qty = parseInt(newQty);
    if (qty <= 0) {
        await api(`/api/portfolios/${pid}/assets/${aid}`, 'DELETE');
    } else {
        await api(`/api/portfolios/${pid}/assets`, 'POST', { asset_id: aid, quantity: qty });
    }
    await renderPortfoliosView();
}

async function updatePortfolioViewLotSize(pid, aid, newLot) {
    const lot = parseInt(newLot);
    if (!isNaN(lot) && lot >= 1) {
        await api(`/api/portfolios/${pid}/assets/${aid}/lot`, 'PUT', { lot_size: lot });
        await renderPortfoliosView();
    }
}

async function deletePortfolioAsset(pid, aid) {
    await api(`/api/portfolios/${pid}/assets/${aid}`, 'DELETE');
    await renderPortfoliosView();
}

async function renderAssetsView() {
    const allAssets = await api('/api/assets') || [];
    allAssets.sort((a, b) => a.name.localeCompare(b.name));

    content.innerHTML = `
        <div class="card" style="margin-bottom: 24px;">
            <input type="text" id="assets-search-filter" class="search-bar" placeholder="Search system asset catalog by ticker name..." oninput="onAssetCatalogSearch()">
        </div>
        <div class="data-table-container">
            <table>
                <thead>
                    <tr>
                        <th>Asset Ticker</th>
                        <th>Latest Parsed Valuation Price</th>
                    </tr>
                </thead>
                <tbody id="assets-catalog-tbody">
                    ${allAssets.map(a => `
                        <tr>
                            <td><strong>${a.name}</strong></td>
                            <td>₽ ${a.price.toFixed(2)}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>`;
}

async function onAssetCatalogSearch() {
    const searchVal = document.getElementById('assets-search-filter').value.toUpperCase();
    const allAssets = await api('/api/assets') || [];
    const filtered = allAssets.filter(a => a.name.toUpperCase().includes(searchVal));
    const tbody = document.getElementById('assets-catalog-tbody');
    if (tbody) {
        tbody.innerHTML = filtered.map(a => `
            <tr>
                <td><strong>${a.name}</strong></td>
                <td>₽ ${a.price.toFixed(2)}</td>
            </tr>
        `).join('');
    }
}

initializeApp();
