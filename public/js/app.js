$(document).ready(function () {
    const API_URL = 'http://localhost:3000/api';
    let cart = [];
    let itemsCache = [];

    const formatRp = (num) => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(num);

    function apiRequest(endpoint, method = 'GET', data = null) {
        const token = sessionStorage.getItem('jwt_token');

        return $.ajax({
            url: API_URL + endpoint,
            method: method,
            contentType: 'application/json',
            data: data ? JSON.stringify(data) : null,
            headers: { 'Authorization': token ? `Bearer ${token}` : '' },
        }).fail(function (xhr) {
            let msg = "Terjadi kesalahan sistem.";
            if (xhr.responseJSON && xhr.responseJSON.error) {
                msg = xhr.responseJSON.error;
            }

            if (xhr.status === 401 && endpoint !== '/login') {
                Swal.fire({
                    icon: 'warning',
                    title: 'Sesi Habis',
                    text: 'Silakan login kembali.',
                    timer: 1500,
                    showConfirmButton: false
                }).then(() => {
                    doLogout();
                });
                return;
            }

            if (xhr.status === 0) {
                msg = "Gagal terhubung ke Backend (Server Mati).";
            }

            Swal.fire({
                icon: 'error',
                title: 'Gagal',
                text: msg,
                confirmButtonColor: '#2563eb'
            });
        });
    }

    function checkAuth() {
        const token = sessionStorage.getItem('jwt_token');
        const user = JSON.parse(sessionStorage.getItem('user'));

        $('#auth-container, #view-dashboard').addClass('hidden-force');
        $('#nav-user').addClass('hidden-force');

        if (token && user) {
            $('#view-dashboard').removeClass('hidden-force');
            $('#nav-user').removeClass('hidden-force');
            $('#user-display').text(user.username);
            loadMasterData();
        } else {
            $('#auth-container').removeClass('hidden-force');
            $('#view-login').removeClass('hidden-force');
        }


        setTimeout(() => {
            $('#app-preloader').fadeOut(300, function () {
                $(this).remove();
            });
        }, 300);
    }

    function doLogout() {
        sessionStorage.clear();
        cart = [];
        window.location.reload();
    }

    $('#btn-logout').click(doLogout);
    $('#link-to-register').click(e => { e.preventDefault(); $('#view-login').addClass('hidden-force'); $('#view-register').removeClass('hidden-force'); });
    $('#link-to-login').click(e => { e.preventDefault(); $('#view-register').addClass('hidden-force'); $('#view-login').removeClass('hidden-force'); });

    $('#form-login').submit(function (e) {
        e.preventDefault();
        const btn = $(this).find('button');
        const originalText = btn.text();

        btn.prop('disabled', true).text('Memproses...');

        apiRequest('/login', 'POST', {
            username: $('#login-username').val(),
            password: $('#login-password').val()
        })
            .done(resp => {
                sessionStorage.setItem('jwt_token', resp.token);
                sessionStorage.setItem('user', JSON.stringify(resp.user));

                Swal.fire({
                    icon: 'success',
                    title: 'Login Berhasil',
                    text: 'Mengalihkan ke dashboard...',
                    timer: 1000,
                    showConfirmButton: false
                }).then(() => {
                    window.location.reload();
                });
            })
            .always(() => {
                btn.prop('disabled', false).text(originalText);
            });
    });

    $('#form-register').submit(function (e) {
        e.preventDefault();
        const btn = $(this).find('button');
        const originalText = btn.text();
        btn.prop('disabled', true).text('Memproses...');

        apiRequest('/register', 'POST', {
            username: $('#reg-username').val(),
            password: $('#reg-password').val()
        })
            .done(() => {
                Swal.fire('Sukses', 'Akun dibuat. Silakan login.', 'success');
                $('#link-to-login').click();
            })
            .always(() => btn.prop('disabled', false).text(originalText));
    });

    $('#mobile-inv-toggle').click(function () {
        if ($(window).width() < 1024) {
            $('#inventory-list-wrapper').slideToggle(200);
            $('#mobile-inv-chevron').toggleClass('rotate-180');
        }
    });


    function loadMasterData() {
        apiRequest('/suppliers').done(function (suppliers) {
            const $select = $('#input-supplier').html('<option value="">-- Pilih Supplier --</option>');
            suppliers.forEach(supp => $select.append(`<option value="${supp.id}">${supp.name}</option>`));
        });

        $('#list-items').html('<li class="text-center py-4"><i class="fas fa-spinner fa-spin"></i> Memuat Data...</li>');

        apiRequest('/items').done(function (items) {
            itemsCache = items;
            renderSidebarList(items);
        });
    }

    function renderSidebarList(items) {
        const $list = $('#list-items');

        if (items.length === 0) {
            $list.html('<li class="text-center py-4 text-slate-400">Tidak ada barang ditemukan</li>');
            $('#item-count').text(0);
            return;
        }

        let htmlContent = items.map(item => {
            let badgeClass = item.stock < 10 ? 'bg-red-100 text-red-600' : 'bg-green-100 text-green-600';
            return `
                <li class="item-row p-3 bg-white hover:bg-blue-50 border border-slate-100 rounded-lg cursor-pointer transition flex justify-between items-center group mb-2" onclick="selectItemFromSidebar(${item.id})">
                    <div class="overflow-hidden">
                        <div class="font-bold text-slate-700 text-sm truncate group-hover:text-primary">${item.name}</div>
                        <div class="text-xs text-slate-400">${formatRp(item.price)}</div>
                    </div>
                    <span class="text-[10px] font-bold px-2 py-1 rounded-full ${badgeClass} flex-none ml-2">${item.stock}</span>
                </li>
            `;
        }).join('');

        $list.html(htmlContent);
        $('#item-count').text(items.length);
        $('#mobile-badge-count').text(items.length);
        updateDropdownOptions(items);
    }

    function updateDropdownOptions(items) {
        const limitedItems = items.slice(0, 100);
        const $select = $('#input-item').html('<option value="">-- Klik barang di Sidebar --</option>');
        limitedItems.forEach(item => {
            $select.append(`<option value="${item.id}">${item.name}</option>`);
        });
        if (items.length > 100) $select.append(`<option value="" disabled>...Gunakan search...</option>`);
    }

    function handleSearch(keyword) {
        const filtered = itemsCache.filter(item => item.name.toLowerCase().includes(keyword));
        renderSidebarList(filtered);
    }

    $('#input-search').on('input', function () { handleSearch($(this).val().toLowerCase()); });
    $('#input-search-mobile').on('input', function () { handleSearch($(this).val().toLowerCase()); });
    $('#btn-refresh').click(loadMasterData);

    window.selectItemFromSidebar = function (id) {
        const item = itemsCache.find(i => i.id == id);
        if (!item) return;

        if ($(`#input-item option[value='${id}']`).length === 0) {
            $('#input-item').append(`<option value="${item.id}">${item.name}</option>`);
        }
        $('#input-item').val(id).trigger('change');

        if ($(window).width() < 1024) {
            $('#inventory-list-wrapper').slideUp(200);
            $('#mobile-inv-chevron').removeClass('rotate-180');
            $('html, body').animate({ scrollTop: $("#input-item").offset().top - 150 }, 500);
        }
        $('#input-qty').focus().select();
    };

    $('#input-item').change(function () {
        const id = $(this).val();
        if (id) {
            const item = itemsCache.find(i => i.id == id);
            $('#info-name').text(item.name);
            $('#info-stock').text(item.stock);
            $('#info-price').text(formatRp(item.price));
            $('#info-box').removeClass('hidden-force');
        } else {
            $('#info-box').addClass('hidden-force');
        }
    });


    $('#btn-add-cart').click(function () {
        const supplierId = $('#input-supplier').val();
        const itemId = $('#input-item').val();
        const qty = parseInt($('#input-qty').val());

        if (!supplierId) return Swal.fire({ icon: 'warning', title: 'Pilih Supplier', text: 'Tentukan supplier dulu.' });
        if (!itemId) return Swal.fire({ icon: 'warning', title: 'Pilih Barang' });
        if (!qty || qty < 1) return Swal.fire({ icon: 'warning', title: 'Qty Salah' });

        const item = itemsCache.find(i => i.id == itemId);
        const existing = cart.find(c => c.item_id == item.id);

        if (existing) existing.qty += qty;
        else cart.push({ item_id: item.id, name: item.name, price: item.price, qty: qty });

        renderCart();
        $('#input-qty').val(1);
    });

    $('#table-cart-body').on('click', '.btn-remove', function () {
        cart.splice($(this).data('index'), 1);
        renderCart();
    });

    $('#btn-clear-cart, #btn-clear-cart-mobile, #btn-clear-cart-sm').click(function () {
        if (cart.length === 0) return;
        Swal.fire({
            title: 'Reset?',
            text: "Keranjang akan dikosongkan.",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#dc2626',
            confirmButtonText: 'Ya, Reset'
        }).then((result) => {
            if (result.isConfirmed) {
                cart = [];
                renderCart();
            }
        });
    });

    function renderCart() {
        const $tbody = $('#table-cart-body').empty();
        let grandTotal = 0;
        const hasItems = cart.length > 0;

        $('#input-supplier').prop('disabled', hasItems);
        hasItems ? $('#supplier-lock-msg').removeClass('hidden-force') : $('#supplier-lock-msg').addClass('hidden-force');
        hasItems ? $('#btn-clear-cart-sm').removeClass('hidden-force') : $('#btn-clear-cart-sm').addClass('hidden-force');

        const btns = $('#btn-submit-order, #btn-clear-cart, #btn-clear-cart-mobile');
        if (hasItems) {
            btns.prop('disabled', false).removeClass('opacity-50 cursor-not-allowed');
            $('#table-cart-foot').removeClass('hidden-force');
        } else {
            btns.prop('disabled', true).addClass('opacity-50 cursor-not-allowed');
            $('#table-cart-foot').addClass('hidden-force');
            $tbody.html(`<tr><td colspan="4" class="px-4 py-8 text-center text-slate-400 italic">Keranjang kosong</td></tr>`);
            return;
        }

        cart.forEach((c, index) => {
            const subtotal = c.price * c.qty;
            grandTotal += subtotal;
            $tbody.append(`
                <tr class="hover:bg-slate-50 transition border-b border-slate-50">
                    <td class="px-4 py-3 align-middle">
                        <div class="text-slate-700 font-medium">${c.name}</div>
                        <div class="text-xs text-slate-400 lg:hidden">${formatRp(c.price)} x ${c.qty}</div>
                    </td>
                    <td class="px-2 py-3 text-center align-middle hidden lg:table-cell">
                        <span class="bg-slate-100 px-2 py-1 rounded text-xs font-bold text-slate-600">${c.qty}</span>
                    </td>
                    <td class="px-4 py-3 text-right font-mono text-xs font-bold text-slate-700 align-middle">
                        ${formatRp(subtotal)}
                    </td>
                    <td class="px-4 py-3 text-center align-middle">
                        <button class="btn-remove text-red-400 hover:text-red-600 w-8 h-8 flex items-center justify-center rounded-full hover:bg-red-50 transition" data-index="${index}">
                            <i class="fas fa-trash"></i>
                        </button>
                    </td>
                </tr>
            `);
        });

        $('#label-grand-total').text(formatRp(grandTotal));
        const tableContainer = $('#table-cart-body').closest('.overflow-y-auto');
        if (tableContainer.length) tableContainer.scrollTop(tableContainer[0].scrollHeight);
    }

    $('#btn-submit-order').click(function () {
        const supplierId = $('#input-supplier').val();
        const btn = $(this);
        const originalContent = btn.html();

        btn.prop('disabled', true).html('<i class="fas fa-spinner fa-spin"></i> Mengirim...');

        apiRequest('/purchase', 'POST', {
            supplier_id: parseInt(supplierId),
            items: cart.map(c => ({ item_id: c.item_id, qty: c.qty }))
        }).done(function () {
            Swal.fire({ icon: 'success', title: 'Order Sukses', timer: 1500 });
            cart = [];
            renderCart();
            $('#input-supplier').val('');
            loadMasterData();
        }).always(() => {
            if (cart.length > 0) btn.prop('disabled', false).html(originalContent);
            else btn.html(originalContent);
        });
    });

    checkAuth();
});