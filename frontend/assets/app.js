const API_BASE = 'http://localhost:8080/api';

// Mostrar/esconder abas
function showTab(tabName) {
    // Esconder todas as abas
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.style.display = 'none';
    });

    // Mostrar aba selecionada
    document.getElementById(tabName).style.display = 'block';

    // Atualizar sidebar
    document.querySelectorAll('.list-group-item').forEach(item => {
        item.classList.remove('active');
    });
    event.target.classList.add('active');

    // Carregar dados conforme necessário
    if (tabName === 'dashboard') {
        carregarStats();
    } else if (tabName === 'leads') {
        carregarLeads();
    } else if (tabName === 'campaigns') {
        carregarCampanhas();
    }
}

// DASHBOARD - Carregar estatísticas
async function carregarStats() {
    try {
        const response = await fetch(`${API_BASE}/stats`);
        const stats = await response.json();

        document.getElementById('statTotalLeads').textContent = stats.total_leads || 0;
        document.getElementById('statLeadsAltos').textContent = stats.leads_altos || 0;
        document.getElementById('statLeadsMedios').textContent = stats.leads_medios || 0;
        document.getElementById('statLeadsBaixos').textContent = stats.leads_baixos || 0;
        document.getElementById('statCampanhas').textContent = stats.total_campanhas || 0;
        document.getElementById('statTaxaConversao').textContent = 
            (stats.taxa_conversao || 0).toFixed(1) + '%';
    } catch (error) {
        console.error('Erro ao carregar stats:', error);
    }
}

// BUSCA - Executar nova busca
async function executarBusca(event) {
    event.preventDefault();

    const termo = document.getElementById('termoBusca').value;
    const categoria = document.getElementById('categoriaBusca').value;

    try {
        const response = await fetch(`${API_BASE}/search`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ termo, categoria })
        });

        const resultado = await response.json();

        const div = document.getElementById('resultadoBusca');
        div.innerHTML = `
            <strong>✅ Busca iniciada!</strong><br>
            ID: ${resultado.id}<br>
            Termo: ${resultado.termo}<br>
            Status: ${resultado.status}
        `;
        div.style.display = 'block';

        document.getElementById('formBusca').reset();

        setTimeout(() => carregarLeads(), 2000);
    } catch (error) {
        alert('Erro ao executar busca: ' + error);
    }
}

// LEADS - Carregar lista de leads
async function carregarLeads() {
    try {
        const nome = document.getElementById('filtroNome')?.value || '';
        const status = document.getElementById('filtroStatus')?.value || '';
        const qualidade = document.getElementById('filtroQualidade')?.value || '';

        let url = `${API_BASE}/leads?limit=100`;
        if (nome) url += `&busca=${encodeURIComponent(nome)}`;
        if (status) url += `&status=${status}`;
        if (qualidade) url += `&qualidade=${qualidade}`;

        const response = await fetch(url);
        const leads = await response.json();

        const tbody = document.getElementById('bodyLeads');
        tbody.innerHTML = '';

        if (!leads || leads.length === 0) {
            tbody.innerHTML = '<tr><td colspan="8" class="text-center text-muted">Nenhum lead encontrado</td></tr>';
            return;
        }

        leads.forEach(lead => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td><strong>${lead.nome}</strong></td>
                <td><a href="tel:${lead.telefone}">${lead.telefone}</a></td>
                <td><a href="mailto:${lead.email}">${lead.email || '-'}</a></td>
                <td>${lead.categoria || '-'}</td>
                <td>
                    <span class="badge ${getQualidadeBadge(lead.qualidade)}">
                        ${lead.qualidade}
                    </span>
                </td>
                <td>
                    <select class="form-select form-select-sm" onchange="atualizarLead('${lead.id}', 'status', this.value)">
                        <option value="novo" ${lead.status === 'novo' ? 'selected' : ''}>Novo</option>
                        <option value="contato" ${lead.status === 'contato' ? 'selected' : ''}>Contato</option>
                        <option value="conversao" ${lead.status === 'conversao' ? 'selected' : ''}>Conversão</option>
                        <option value="perdido" ${lead.status === 'perdido' ? 'selected' : ''}>Perdido</option>
                    </select>
                </td>
                <td>${lead.perfil || '-'}</td>
                <td>
                    <button class="btn btn-sm btn-warning" onclick="editarLead('${lead.id}')">✏️</button>
                    <button class="btn btn-sm btn-danger" onclick="deletarLead('${lead.id}')">🗑️</button>
                </td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        console.error('Erro ao carregar leads:', error);
    }
}

// LEADS - Atualizar status
async function atualizarLead(id, campo, valor) {
    try {
        // Primeiro, fazer GET para pegar dados atuais
        const response1 = await fetch(`${API_BASE}/leads/${id}`);
        const lead = await response1.json();

        // Atualizar o campo
        lead[campo] = valor;

        const response = await fetch(`${API_BASE}/leads/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(lead)
        });

        if (response.ok) {
            carregarLeads();
        }
    } catch (error) {
        alert('Erro ao atualizar lead: ' + error);
    }
}

// LEADS - Deletar
async function deletarLead(id) {
    if (!confirm('Tem certeza que deseja deletar este lead?')) return;

    try {
        const response = await fetch(`${API_BASE}/leads/${id}`, { method: 'DELETE' });
        if (response.ok) {
            carregarLeads();
        }
    } catch (error) {
        alert('Erro ao deletar: ' + error);
    }
}

// CAMPANHAS - Carregar lista
async function carregarCampanhas() {
    try {
        const response = await fetch(`${API_BASE}/campaigns`);
        const campanhas = await response.json();

        const tbody = document.getElementById('bodyCampanhas');
        tbody.innerHTML = '';

        if (!campanhas || campanhas.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="text-center text-muted">Nenhuma campanha criada</td></tr>';
            return;
        }

        campanhas.forEach(c => {
            const taxa = c.contatos_feitos > 0 
                ? ((c.conversoes / c.contatos_feitos) * 100).toFixed(1) + '%'
                : '0%';

            const row = document.createElement('tr');
            row.innerHTML = `
                <td><strong>${c.nome}</strong></td>
                <td>${c.tipo}</td>
                <td>${c.total_leads}</td>
                <td>${c.contatos_feitos}</td>
                <td>${c.conversoes}</td>
                <td>${taxa}</td>
            `;
            tbody.appendChild(row);
        });
    } catch (error) {
        console.error('Erro ao carregar campanhas:', error);
    }
}

// CAMPANHAS - Abrir form
function abrirFormCampanha() {
    const div = document.getElementById('formCampanhaDiv');
    div.style.display = div.style.display === 'none' ? 'block' : 'none';
}

// CAMPANHAS - Salvar
async function salvarCampanha(event) {
    event.preventDefault();

    const nome = document.getElementById('nomeCampanha').value;
    const tipo = document.getElementById('tipoCampanha').value;

    try {
        const response = await fetch(`${API_BASE}/campaigns`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ nome, tipo })
        });

        if (response.ok) {
            alert('✅ Campanha criada com sucesso!');
            document.getElementById('formCampanhaDiv').style.display = 'none';
            carregarCampanhas();
        }
    } catch (error) {
        alert('Erro ao criar campanha: ' + error);
    }
}

// EXPORT - CSV
async function exportarCSV() {
    window.location.href = `${API_BASE}/export/csv`;
}

// EXPORT - JSON
async function exportarJSON() {
    window.location.href = `${API_BASE}/export/json`;
}

// HELPERS
function getQualidadeBadge(qualidade) {
    const badges = {
        'high': 'bg-success',
        'medium': 'bg-warning',
        'low': 'bg-danger'
    };
    return badges[qualidade] || 'bg-secondary';
}

function editarLead(id) {
    alert('Edição individual - implementar modal');
}

// Carregar stats ao iniciar
document.addEventListener('DOMContentLoaded', () => {
    carregarStats();
});
