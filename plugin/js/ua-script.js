// Shared JavaScript for unRAID Agent plugin pages.
// Loaded once from the parent tab page via <script src> so it is always
// available for all AJAX-loaded child tabs (Unraid's tab system may not
// execute <script> blocks in AJAX responses). All handlers use event
// delegation (bind on document) so they work regardless of when content
// is injected.

if (typeof jQuery === 'undefined') {
    document.addEventListener('DOMContentLoaded', function() {
        if (typeof jQuery !== 'undefined') jQuery(window).trigger('ua:jquery-ready');
    });
}

jQuery(function($) {
    'use strict';

    var initialized = false;

    function init() {
        if (initialized) return;
        initialized = true;
        restoreActiveTab();
        expandFirstDomain();
    }

    // ---- Sub-tab switching (shared between Permissions + Content pages) ----

    var PERMS_KEY   = 'unraid-agent-perms-tab';
    var CONTENT_KEY = 'unraid-agent-content-tab';
    var TOOLSET_KEY = 'unraid-agent-toolset-tab';

    $(document).on('click', '.ua-subtab[data-uatab]', function() {
        var t = $(this).data('uatab');
        if (!t) return;
        var isToolGroup = (t.indexOf('tg-') === 0);

        if (isToolGroup) {
            $('.ua-subtab[data-uatab^="tg-"]').removeClass('active');
            $('.ua-pane[data-uatab^="tg-"], .ac-pane[data-uatab^="tg-"]').removeClass('active');
            $(this).addClass('active');
            $('.ua-pane[data-uatab="' + t + '"], .ac-pane[data-uatab="' + t + '"]').addClass('active');
            localStorage.setItem(TOOLSET_KEY, t);
        } else {
            $('.ua-subtab[data-uatab]:not([data-uatab^="tg-"])').removeClass('active');
            $('.ua-pane[data-uatab]:not([data-uatab^="tg-"]), .ac-pane[data-uatab]:not([data-uatab^="tg-"])').removeClass('active');
            $(this).addClass('active');
            $('.ua-pane[data-uatab="' + t + '"], .ac-pane[data-uatab="' + t + '"]').addClass('active');

            if ($('.ua-domain').length) {
                localStorage.setItem(PERMS_KEY, t);
                if (t === 'tools') {
                    var savedTg = localStorage.getItem(TOOLSET_KEY);
                    if (!savedTg || !$('.ua-subtab[data-uatab="' + savedTg + '"]').length) {
                        savedTg = 'tg-core';
                    }
                    $('.ua-subtab[data-uatab^="tg-"]').removeClass('active');
                    $('.ua-pane[data-uatab^="tg-"], .ac-pane[data-uatab^="tg-"]').removeClass('active');
                    $('.ua-subtab[data-uatab="' + savedTg + '"]').addClass('active');
                    $('.ua-pane[data-uatab="' + savedTg + '"], .ac-pane[data-uatab="' + savedTg + '"]').addClass('active');
                    localStorage.setItem(TOOLSET_KEY, savedTg);
                }
            } else {
                localStorage.setItem(CONTENT_KEY, t);
            }
        }
    });

    function restoreActiveTab() {
        var $panes = $('.ua-pane[data-uatab], .ac-pane[data-uatab]');
        if ($panes.length === 0) return;
        var $activePane = $('.ua-pane.active, .ac-pane.active');
        if ($activePane.length) return;

        var key = $('.ua-domain').length ? PERMS_KEY : CONTENT_KEY;
        var saved = localStorage.getItem(key);
        var $target = null;
        if (saved) {
            $target = $('.ua-subtab[data-uatab="' + saved + '"]');
        }
        if (!$target || !$target.length) {
            $target = $('.ua-subtab[data-uatab]:not([data-uatab^="tg-"]):first');
        }
        if ($target && $target.length) {
            $target.click();
        } else if ($('.ua-subtab[data-uatab^="tg-"]:first').length) {
            $('.ua-subtab[data-uatab^="tg-"]:first').click();
        }
    }

    function expandFirstDomain() {
        var $firstHead = $('.ua-domain-head:first');
        var $firstTools = $firstHead.next('.ua-domain-tools');
        if ($firstHead.length && !$firstTools.is(':visible')) {
            $firstHead.click();
        }
    }

    // Run initialization on DOM ready and when content might be present
    init();

    // Also listen for a custom event in case jQuery loads after DOM ready
    $(document).on('ua:jquery-ready', init);

    // ---- Tool domain expand/collapse (Permissions page only) ----
    $(document).on('click', '.ua-domain-head', function() {
        var $tools = $(this).next('.ua-domain-tools');
        $tools.toggle();
        $(this).find('.ua-domain-caret').text($tools.is(':visible') ? '\u25BE' : '\u25B8');
    });

    // ---- Modal handling (shared between Permissions + Content pages) ----

    var currentModal = null;

    $(document).on('click', '.ua-card[data-entity]', function(e) {
        if ($(e.target).hasClass('ac-export') || $(e.target).closest('.ac-export').length) return;
        e.preventDefault();
        e.stopPropagation();

        var card = $(this);
        var entity = card.data('entity');
        if (!entity) return;

        var type = card.data('type');
        var icon = card.data('icon');
        var perms = card.data('perms') || {};
        var isSelf = (type === 'plugins' && entity === 'unraid-agent');

        $('#ua-modal-name').text(entity);
        if (icon) {
            $('#ua-modal-icon').html('<img src="' + icon + '" style="width:44px;height:44px;object-fit:contain;">');
        } else {
            var iconClass = type === 'vms' ? 'fa-desktop' : (type === 'plugins' ? 'fa-plug' : 'fa-cubes');
            var iconColor = type === 'vms' ? '#9b59b6' : (type === 'plugins' ? '#2ecc71' : '#0db7ed');
            $('#ua-modal-icon').html('<i class="fa ' + iconClass + '" style="font-size:40px;color:' + iconColor + ';"></i>');
        }

        currentModal = { type: type, entity: entity, icon: icon, perms: perms };

        var ACTIONS = {
            containers: [['start','Start'],['stop','Stop'],['restart','Restart'],['pause','Pause'],['unpause','Unpause'],['remove','Remove'],['update','Update']],
            vms: [['start','Start'],['stop','Stop'],['pause','Pause'],['resume','Resume'],['force_stop','Force Stop'],['reboot','Reboot'],['reset','Reset']],
            plugins: [['remove','Remove']]
        };

        if (isSelf) {
            $('#ua-modal-rows').html(
                '<div style="padding:12px 4px;line-height:1.6;">' +
                '<b>Protected plugin.</b><br>' +
                'unraid-agent can never be removed through the MCP server — that would shut down ' +
                'this very channel. If removal is ever intended, use the Unraid WebGUI Plugins page.' +
                '</div>');
            $('#ua-modal-save').hide();
            $('#ua-modal-clear').hide();
        } else {
            var rows = '';
            $.each(ACTIONS[type], function(i, pair) {
                var action = pair[0], label = pair[1];
                var val = perms[action] || 'inherit';
                rows += '<div class="ua-perm-row"><div class="ua-perm-pair">' +
                    '<span class="ua-perm-label">' + label + '</span>' +
                    '<select class="ua-perm-select" data-action="' + action + '">' +
                    '<option value="inherit"' + (val === 'inherit' ? ' selected' : '') + '>Use Global</option>' +
                    '<option value="allow"' + (val === 'allow' ? ' selected' : '') + '>Allow</option>' +
                    '<option value="deny"' + (val === 'deny' ? ' selected' : '') + '>Deny</option>' +
                    '</select></div></div>';
            });
            $('#ua-modal-rows').html(rows);
            $('#ua-modal-save').show();
            $('#ua-modal-clear').show();
        }
        $('#ua-modal-err').text('');
        $('#ua-modal').fadeIn(120);
    });

    function closeModal() {
        $('#ua-modal').fadeOut(120);
        currentModal = null;
    }

    // Clicking outside modals closes them
    $(document).on('click', '#ac-modal', function(e) {
        if (e.target === this) { $('#ac-modal').fadeOut(120); currentModal = null; }
    });
    $(document).on('click', '#ua-modal', function(e) {
        if (e.target === this) closeModal();
    });

    $(document).on('click', '.ac-skill-card, .ac-memory-card', function(e) {
        if ($(e.target).hasClass('ac-export') || $(e.target).closest('.ac-export').length) return;
        var card = $(this);
        var file = card.data('file');
        if (!file) return;

        var kind = card.hasClass('ac-skill-card') ? 'skill' : 'memory';
        var title = card.data('name');
        currentModal = { kind: kind, data: { name: title, file: file, source: card.data('source'), scope: card.data('scope'), isNew: false } };

        loadFile(file, function(content) {
            openContentModal(kind, title, content);
        });
    });

    function loadFile(file, cb) {
        $.post('/plugins/unraid-agent/php/read-content.php',
            { csrf_token: window.UA_CSRF || '', file: file },
            function(resp) {
                cb(resp && resp.ok ? resp.content : '', null);
            }, 'json').fail(function(xhr) {
                var msg = '// Failed to load content (HTTP ' + xhr.status + ')';
                try { var r = JSON.parse(xhr.responseText); if (r.error) msg += ' — ' + r.error; } catch(e) {}
                cb(msg, msg);
            });
    }

    function openContentModal(kind, title, content) {
        var isSkill = (kind === 'skill');
        var data = currentModal.data;
        $('#ac-modal-title').text(data.isNew ? (isSkill ? 'New Skill' : 'New Memory') : (isSkill ? 'Edit Skill' : 'Edit Memory'));
        $('#ac-modal-icon').attr('class', 'fa ' + (isSkill ? 'fa-graduation-cap' : 'fa-sticky-note'))
            .css({ fontSize: '34px', color: isSkill ? '#ff8c2f' : '#4a90d9' });

        var desc = isSkill ? data.description || '' : '';
        if (!data.isNew && isSkill && content) {
            var m;
            if ((m = content.match(/^---\s*\n(.*?)\n---\s*\n?(.*)$/s))) {
                var front = m[1], body = m[2];
                var nm = front.match(/^name:\s*(.+)$/m);
                var dm = front.match(/^description:\s*(.+)$/m/);
                if (nm) data.name = nm[1].trim();
                if (dm) desc = dm[1].trim();
                content = body;
            }
        }

        $('#ac-name').val(data.name || '').prop('disabled', !data.isNew && !isSkill);
        $('#ac-desc').val(desc);
        $('#ac-field-desc').toggle(isSkill);
        $('#ac-field-name label').text(isSkill
            ? 'Name (lowercase letters, digits, hyphens)'
            : 'Name (lowercase letters, digits, hyphens) — scope: ' + (data.scope || 'default'));
        $('#ac-content').val(content || '');
        $('#ac-delete').toggle(!data.isNew);
        $('#ac-err').text('');
        $('#ac-modal').fadeIn(120);
    }

    // ---- Button click handlers ----

    $(document).on('click', '#ac-cancel', function(e) {
        e.preventDefault();
        $('#ac-modal').fadeOut(120);
        currentModal = null;
    });

    $(document).on('click', '#ua-modal-cancel', function(e) {
        e.preventDefault();
        closeModal();
    });

    $(document).on('click', '#ac-save', function(e) {
        e.preventDefault();
        if (!currentModal) return;
        var isSkill = (currentModal.kind === 'skill');
        var payload = {
            csrf_token: window.UA_CSRF || '',
            name: $('#ac-name').val(),
            content: $('#ac-content').val()
        };
        if (isSkill) payload.description = $('#ac-desc').val();
        if (!isSkill) payload.scope = currentModal.data.scope || 'custom';
        $('#ac-err').text('Saving…');
        $.post('/plugins/unraid-agent/php/save-content.php', $.extend(payload, { kind: currentModal.kind }), function(resp) {
            if (resp && resp.ok) { location.reload(); }
            else { $('#ac-err').text((resp && resp.error) ? resp.error : 'Save failed'); }
        }, 'json').fail(function(xhr) {
            var msg = 'Save failed (' + xhr.status + ')';
            try { var r = JSON.parse(xhr.responseText); if (r.error) msg = r.error; } catch(err) {}
            $('#ac-err').text(msg);
        });
    });

    $(document).on('click', '#ac-delete', function(e) {
        e.preventDefault();
        if (!currentModal || !currentModal.data) return;
        if (!confirm('Delete this entry?')) return;
        $('#ac-err').text('Deleting…');
        $.post('/plugins/unraid-agent/php/save-content.php', {
            csrf_token: window.UA_CSRF || '',
            kind: currentModal.kind,
            action: 'delete',
            file: currentModal.data.file
        }, function(resp) {
            if (resp && resp.ok) { location.reload(); }
            else { $('#ac-err').text((resp && resp.error) ? resp.error : 'Delete failed'); }
        }, 'json').fail(function(xhr) {
            var msg = 'Delete failed (' + xhr.status + ')';
            try { var r = JSON.parse(xhr.responseText); if (r.error) msg = r.error; } catch(err) {}
            $('#ac-err').text(msg);
        });
    });

    // ---- Permission modal save/clear ----

    $(document).on('click', '#ua-modal-save', function(e) {
        e.preventDefault();
        if (!currentModal) return;
        $('#ua-modal-err').text('Saving…');
        var perms = {};
        $('.ua-perm-select').each(function() {
            perms[$(this).data('action')] = $(this).val();
        });
        $.ajax({
            url: '/plugins/unraid-agent/php/save-perms.php',
            method: 'POST',
            data: {
                csrf_token: window.UA_CSRF || '',
                entity_type: currentModal.type,
                entity: currentModal.entity,
                perms: JSON.stringify(perms)
            },
            success: function(resp) {
                if (resp && resp.ok) { location.reload(); }
                else { $('#ua-modal-err').text((resp && resp.error) ? resp.error : 'Save failed'); }
            },
            error: function(xhr) {
                var msg = 'Save failed (' + xhr.status + ')';
                try { var r = JSON.parse(xhr.responseText); if (r.error) msg = r.error; } catch(err) {}
                $('#ua-modal-err').text(msg);
            }
        });
    });

    $(document).on('click', '#ua-modal-clear', function(e) {
        e.preventDefault();
        if (!currentModal) return;
        $('#ua-modal-err').text('Clearing…');
        $.ajax({
            url: '/plugins/unraid-agent/php/save-perms.php',
            method: 'POST',
            data: {
                csrf_token: window.UA_CSRF || '',
                entity_type: currentModal.type,
                entity: currentModal.entity,
                perms: JSON.stringify({})
            },
            success: function(resp) {
                if (resp && resp.ok) { location.reload(); }
                else { $('#ua-modal-err').text((resp && resp.error) ? resp.error : 'Clear failed'); }
            },
            error: function(xhr) {
                var msg = 'Clear failed (' + xhr.status + ')';
                try { var r = JSON.parse(xhr.responseText); if (r.error) msg = r.error; } catch(err) {}
                $('#ua-modal-err').text(msg);
            }
        });
    });

    // ---- Content action buttons ----

    $(document).on('click', '#ac-new-skill', function() {
        currentModal = { kind: 'skill', data: { isNew: true, source: 'custom', name: '', content: '', description: '' } };
        openContentModal('skill', 'New Skill', '');
    });

    $(document).on('click', '#ac-new-memory', function() {
        currentModal = { kind: 'memory', data: { isNew: true, scope: 'custom', name: '', content: '' } };
        openContentModal('memory', 'New Memory', '');
    });

    // ---- Export (download) ----

    $(document).on('click', '.ac-export', function(e) {
        e.stopPropagation();
        var el = $(this);
        var kind = el.data('kind') || 'skill';
        var url = '/plugins/unraid-agent/php/export-content.php?kind=' + encodeURIComponent(kind) + '&name=' + encodeURIComponent(el.data('name'));
        if (kind === 'skill') {
            url += '&source=' + encodeURIComponent(el.data('source') === 'default' ? 'defaults' : 'custom');
        } else {
            url += '&scope=' + encodeURIComponent(el.data('scope'));
        }
        window.location = url;
    });

});
