// JavaScript Frontend para Sistema Contable LATAM
// Interfaz para demostrar R2Lang DSL con Siigo

const API_BASE = '/api';
let regions = [];
let transactions = [];

// Configuración de regiones para display
const regionConfig = {
    'MX': { flag: '🇲🇽', name: 'México', color: 'success' },
    'COL': { flag: '🇨🇴', name: 'Colombia', color: 'info' },
    'AR': { flag: '🇦🇷', name: 'Argentina', color: 'primary' },
    'CH': { flag: '🇨🇱', name: 'Chile', color: 'warning' },
    'UY': { flag: '🇺🇾', name: 'Uruguay', color: 'secondary' },
    'EC': { flag: '🇪🇨', name: 'Ecuador', color: 'dark' },
    'PE': { flag: '🇵🇪', name: 'Perú', color: 'danger' }
};

// Inicialización
document.addEventListener('DOMContentLoaded', function() {
    console.log('🚀 Sistema Contable LATAM - Frontend inicializado');
    loadRegions();
    loadStats();
    
    // Auto-refresh stats cada 30 segundos
    setInterval(loadStats, 30000);
});

// Funciones de API
async function apiCall(endpoint, options = {}) {
    try {
        showLoading(true);
        const response = await fetch(API_BASE + endpoint, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Error en la API');
        }
        
        return data;
    } catch (error) {
        console.error('Error en API:', error);
        showError('Error de conexión: ' + error.message);
        throw error;
    } finally {
        showLoading(false);
    }
}

// Cargar regiones
async function loadRegions() {
    try {
        const response = await apiCall('/regions');
        regions = response.data;
        
        displayRegions();
        populateRegionSelect();
        
        console.log(`✓ ${regions.length} regiones cargadas`);
    } catch (error) {
        console.error('Error cargando regiones:', error);
    }
}

// Cargar estadísticas
async function loadStats() {
    try {
        const response = await apiCall('/stats');
        const stats = response.data;
        
        // Actualizar estadísticas en UI
        document.getElementById('totalRegions').textContent = stats.regions_available || 7;
        document.getElementById('totalTransactions').textContent = stats.transactions.total_transactions || 0;
        document.getElementById('totalAmount').textContent = formatCurrency(stats.transactions.total_amount || 0);
        document.getElementById('systemStatus').textContent = 'ACTIVO';
        
        console.log('✓ Estadísticas actualizadas');
    } catch (error) {
        console.error('Error cargando estadísticas:', error);
    }
}

// Mostrar regiones
function displayRegions() {
    const container = document.getElementById('regionsContainer');
    
    if (regions.length === 0) {
        container.innerHTML = '<div class="col-12 text-center text-muted">No se pudieron cargar las regiones</div>';
        return;
    }
    
    container.innerHTML = regions.map(region => `
        <div class="col-md-6 col-lg-3 mb-3">
            <div class="card region-card border-${regionConfig[region.code]?.color || 'primary'}" 
                 onclick="selectRegion('${region.code}')">
                <div class="card-body text-center">
                    <div class="h2 mb-2">${regionConfig[region.code]?.flag || '🏴'}</div>
                    <h6 class="card-title">${region.name}</h6>
                    <small class="text-muted">
                        ${region.currency} | IVA: ${(region.tax_rate * 100).toFixed(1)}%
                    </small>
                </div>
            </div>
        </div>
    `).join('');
}

// Poblar select de regiones
function populateRegionSelect() {
    const select = document.getElementById('regionSelect');
    select.innerHTML = '<option value="">Seleccionar...</option>' +
        regions.map(region => `
            <option value="${region.code}">
                ${regionConfig[region.code]?.flag || ''} ${region.name} (${region.currency})
            </option>
        `).join('');
}

// Seleccionar región (para futuras funcionalidades)
function selectRegion(code) {
    const region = regions.find(r => r.code === code);
    if (region) {
        document.getElementById('regionSelect').value = code;
        showInfo(`Región seleccionada: ${regionConfig[code]?.flag} ${region.name}`);
    }
}

// Demo completo
async function runCompleteDemo() {
    try {
        showInfo('🚀 Iniciando demo completo para todas las regiones LATAM...');
        
        const response = await apiCall('/demo', { method: 'POST' });
        
        showSuccess(`✓ Demo completado: ${response.transactions_generated} transacciones generadas para ${response.regions_processed} regiones`);
        
        // Cargar transacciones generadas
        loadTransactions();
        loadStats();
        
    } catch (error) {
        console.error('Error en demo completo:', error);
    }
}

// Procesar transacción manual
async function processManualTransaction() {
    const type = document.getElementById('transactionType').value;
    const region = document.getElementById('regionSelect').value;
    const amount = document.getElementById('amountInput').value;
    
    if (!region || !amount) {
        showError('Por favor completa todos los campos');
        return;
    }
    
    if (parseFloat(amount) <= 0) {
        showError('El monto debe ser mayor a cero');
        return;
    }
    
    try {
        const endpoint = type === 'sale' ? '/transactions/sale' : '/transactions/purchase';
        const response = await apiCall(endpoint, {
            method: 'POST',
            body: JSON.stringify({ region, amount: parseFloat(amount) })
        });
        
        showSuccess(`✓ ${type === 'sale' ? 'Venta' : 'Compra'} procesada: ${response.data.transactionId}`);
        
        // Limpiar formulario
        document.getElementById('amountInput').value = '';
        
        // Actualizar datos
        loadTransactions();
        loadStats();
        
        // Mostrar resultado
        displayTransactionResult(response.data);
        
    } catch (error) {
        console.error('Error procesando transacción:', error);
    }
}

// Cargar transacciones
async function loadTransactions() {
    try {
        const response = await apiCall('/transactions');
        transactions = response.data;
        
        displayTransactions();
        
    } catch (error) {
        console.error('Error cargando transacciones:', error);
    }
}

// Mostrar transacciones
function displayTransactions() {
    const container = document.getElementById('resultsContainer');
    
    if (transactions.length === 0) {
        container.innerHTML = `
            <div class="text-center text-muted">
                <i class="fas fa-info-circle fa-3x mb-3"></i>
                <p>No hay transacciones registradas.</p>
            </div>
        `;
        return;
    }
    
    container.innerHTML = transactions.slice(0, 10).map(tx => `
        <div class="card result-card mb-3">
            <div class="card-body">
                <div class="row">
                    <div class="col-md-8">
                        <h6 class="card-title">
                            ${regionConfig[tx.region]?.flag || ''} 
                            ${tx.type === 'sale' ? 'Venta' : 'Compra'} - ${tx.country}
                            <span class="badge bg-success ms-2">${tx.status || 'VALIDADO'}</span>
                        </h6>
                        <p class="card-text">
                            <strong>ID:</strong> ${tx.id}<br>
                            <strong>Monto:</strong> ${formatCurrency(tx.amount)} ${tx.currency}<br>
                            <strong>IVA:</strong> ${formatCurrency(tx.tax)} ${tx.currency}<br>
                            <strong>Total:</strong> ${formatCurrency(tx.total)} ${tx.currency}
                        </p>
                    </div>
                    <div class="col-md-4 text-end">
                        <small class="text-muted">
                            ${new Date(tx.created_at).toLocaleString('es-ES')}<br>
                            Normativa: ${tx.compliance}
                        </small>
                    </div>
                </div>
            </div>
        </div>
    `).join('');
}

// Mostrar resultado de transacción individual
function displayTransactionResult(transaction) {
    const container = document.getElementById('resultsContainer');
    
    const resultHtml = `
        <div class="card result-card mb-3 border-success">
            <div class="card-header bg-success text-white">
                <i class="fas fa-check-circle me-2"></i>Transacción Procesada - ${transaction.transactionId}
            </div>
            <div class="card-body">
                <div class="row">
                    <div class="col-md-6">
                        <h6>${regionConfig[transaction.region]?.flag} ${transaction.country || transaction.region}</h6>
                        <p>
                            <strong>Tipo:</strong> ${transaction.type === 'sale' ? 'Venta' : 'Compra'}<br>
                            <strong>Monto Base:</strong> ${formatCurrency(transaction.amount)} ${transaction.currency}<br>
                            <strong>IVA:</strong> ${formatCurrency(transaction.tax)} ${transaction.currency}<br>
                            <strong>Total:</strong> ${formatCurrency(transaction.total)} ${transaction.currency}
                        </p>
                    </div>
                    <div class="col-md-6">
                        <h6>Asientos Contables:</h6>
                        <div class="code-block">
                            <strong>DEBE:</strong><br>
                            ${transaction.accounts.debit.map(acc => 
                                `${acc.account}: ${formatCurrency(acc.amount)} ${transaction.currency}`
                            ).join('<br>')}<br><br>
                            <strong>HABER:</strong><br>
                            ${transaction.accounts.credit.map(acc => 
                                `${acc.account}: ${formatCurrency(acc.amount)} ${transaction.currency}`
                            ).join('<br>')}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    container.insertAdjacentHTML('afterbegin', resultHtml);
}

// Ejecutar DSL directo
async function executeDSL() {
    const engine = document.getElementById('dslEngine').value;
    const command = document.getElementById('dslCommand').value.trim();
    
    if (!command) {
        showError('Por favor ingresa un comando DSL');
        return;
    }
    
    try {
        const response = await apiCall('/dsl/execute', {
            method: 'POST',
            body: JSON.stringify({ engine, command })
        });
        
        showSuccess(`✓ Comando DSL ejecutado: ${command}`);
        
        // Mostrar resultado
        const resultContainer = document.getElementById('dslResult');
        resultContainer.innerHTML = `
            <div class="card border-info">
                <div class="card-header bg-info text-white">
                    <i class="fas fa-terminal me-2"></i>Resultado DSL - ${engine}
                </div>
                <div class="card-body">
                    <div class="code-block">
                        <strong>Comando:</strong> ${command}<br>
                        <strong>Motor:</strong> ${engine}<br><br>
                        <strong>Resultado:</strong><br>
                        ${JSON.stringify(response.data, null, 2)}
                    </div>
                </div>
            </div>
        `;
        
        // Limpiar comando
        document.getElementById('dslCommand').value = '';
        
        // Actualizar stats si fue una transacción
        if (response.data && response.data.transactionId) {
            loadStats();
        }
        
    } catch (error) {
        console.error('Error ejecutando DSL:', error);
    }
}

// Limpiar resultados
function clearResults() {
    document.getElementById('resultsContainer').innerHTML = `
        <div class="text-center text-muted">
            <i class="fas fa-info-circle fa-3x mb-3"></i>
            <p>No hay transacciones aún. Ejecuta el demo completo o procesa una transacción manual.</p>
        </div>
    `;
    
    document.getElementById('dslResult').innerHTML = '';
    
    showInfo('✓ Resultados limpiados');
}

// Utilidades de UI
function showLoading(show) {
    const loading = document.querySelector('.loading');
    loading.style.display = show ? 'block' : 'none';
}

function showSuccess(message) {
    showToast(message, 'success');
}

function showError(message) {
    showToast(message, 'danger');
}

function showInfo(message) {
    showToast(message, 'info');
}

function showToast(message, type = 'info') {
    // Crear toast dinámico
    const toastHtml = `
        <div class="toast align-items-center text-white bg-${type} border-0" role="alert">
            <div class="d-flex">
                <div class="toast-body">
                    ${message}
                </div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        </div>
    `;
    
    // Crear contenedor si no existe
    let toastContainer = document.getElementById('toastContainer');
    if (!toastContainer) {
        toastContainer = document.createElement('div');
        toastContainer.id = 'toastContainer';
        toastContainer.className = 'toast-container position-fixed top-0 end-0 p-3';
        toastContainer.style.zIndex = '1100';
        document.body.appendChild(toastContainer);
    }
    
    // Agregar toast
    toastContainer.insertAdjacentHTML('beforeend', toastHtml);
    
    // Mostrar toast
    const toastElement = toastContainer.lastElementChild;
    const toast = new bootstrap.Toast(toastElement, { autohide: true, delay: 4000 });
    toast.show();
    
    // Remover elemento después de ocultarse
    toastElement.addEventListener('hidden.bs.toast', function() {
        toastElement.remove();
    });
}

function formatCurrency(amount) {
    if (typeof amount === 'string') {
        amount = parseFloat(amount);
    }
    if (isNaN(amount)) return '$0.00';
    
    return new Intl.NumberFormat('es-ES', {
        style: 'currency',
        currency: 'USD',
        minimumFractionDigits: 2
    }).format(amount).replace('US$', '$');
}

// Funciones de teclado para desarrollo
document.addEventListener('keydown', function(e) {
    // Ctrl+D para demo completo
    if (e.ctrlKey && e.key === 'd') {
        e.preventDefault();
        runCompleteDemo();
    }
    
    // Ctrl+R para recargar stats
    if (e.ctrlKey && e.key === 'r') {
        e.preventDefault();
        loadStats();
        loadRegions();
    }
    
    // Ctrl+L para limpiar
    if (e.ctrlKey && e.key === 'l') {
        e.preventDefault();
        clearResults();
    }
});

// Agregar ejemplos de comandos DSL
const dslExamples = [
    'venta MX 50000',
    'venta COL 75000', 
    'venta AR 100000',
    'compra CH 25000',
    'compra UY 30000',
    'compra EC 15000',
    'compra PE 40000',
    'consultar config MX',
    'consultar config COL'
];

// Click en input DSL para mostrar ejemplos
document.addEventListener('DOMContentLoaded', function() {
    const dslInput = document.getElementById('dslCommand');
    if (dslInput) {
        dslInput.addEventListener('focus', function() {
            if (!this.value) {
                const randomExample = dslExamples[Math.floor(Math.random() * dslExamples.length)];
                this.placeholder = randomExample;
            }
        });
    }
});

console.log('✅ Sistema frontend cargado correctamente');
console.log('🎮 Atajos de teclado:');
console.log('   Ctrl+D: Demo completo');
console.log('   Ctrl+R: Recargar datos'); 
console.log('   Ctrl+L: Limpiar resultados');