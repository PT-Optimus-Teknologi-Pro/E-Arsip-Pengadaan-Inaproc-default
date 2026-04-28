content = open('views/paket/dok-persiapan.html', 'rb').read()
lines = content.split(b'\r\n')

new_lines = []
i = 0
replaced = False

while i < len(lines):
    line_str = lines[i].decode('utf-8', errors='replace')
    
    # Detect the start: inside <td>, first line is <ol>, second is the for loop, third has latestDocId
    if (not replaced 
        and b'<ol>' in lines[i] 
        and i + 2 < len(lines)
        and b'{% for doc in dokPersiapan %}' in lines[i+1]
        and b'latestDocId' in lines[i+2]):
        
        # Find end of this block (</ol>)
        end_idx = i + 1
        while end_idx < len(lines) and b'</ol>' not in lines[end_idx]:
            end_idx += 1
        
        indent = '                    '
        new_block = [
            f'{indent}<ol>',
            f'{indent}{{% set persetujuanShown = false %}}',
            f'{indent}{{% for doc in dokPersiapan %}}',
            f'{indent}    {{% if doc.ChkId == obj.ID and not persetujuanShown %}}',
            f'{indent}        {{% set persetujuanShown = true %}}',
            f'{indent}        {{% for p in doc.Persetujuan() %}}',
            f'{indent}        <li class="mb-1">',
            f'{indent}            <span class="small">{{{{ p.Pegawai().PegNama }}}}</span>',
            f'{indent}            {{% if p.Status %}}',
            f'{indent}            <span class="text-success"><span style="width:16px;" data-feather="check-square"></span></span>',
            f'{indent}            {{% else %}}',
            f'{indent}            <span class="text-muted"><span style="width:16px;" data-feather="square"></span></span>',
            f'{indent}            {{% if (p.PegId == id or isAdmin) and not isLocked %}}',
            f"{indent}            <a href=\"javascript:void(0)\" onclick=\"handleSetuju('{{{{ doc.ID }}}}')\" class=\"btn btn-sm btn-primary ms-1 py-0 px-2\" style=\"font-size: 0.7rem;\">Persetujuan</a>",
            f'{indent}            {{% endif %}}',
            f'{indent}            {{% endif %}}',
            f'{indent}        </li>',
            f'{indent}        {{% endfor %}}',
            f'{indent}    {{% endif %}}',
            f'{indent}{{% endfor %}}',
            f'{indent}</ol>',
        ]
        
        for nl in new_block:
            new_lines.append(nl.encode('utf-8'))
        
        i = end_idx + 1  # skip past </ol>
        replaced = True
        continue
    
    new_lines.append(lines[i])
    i += 1

if replaced:
    new_content = b'\r\n'.join(new_lines)
    open('views/paket/dok-persiapan.html', 'wb').write(new_content)
    print('SUCCESS')
else:
    print('NOT FOUND')
    for idx, l in enumerate(lines[150:175], start=150):
        print(str(idx) + ': ' + repr(l[:80]))
