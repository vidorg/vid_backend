import os
import sys
import re
import yaml
import json


def parsePlain(content) -> {}:
    tokens = re.findall(r'// @(.+)', content)
    d = {}
    for token in tokens:
        sp = re.split(r'[ \t]', token)
        val = ' '.join(sp[1:])
        sp[0] = sp[0].lower()
        if '.' in sp[0]:
            spd = sp[0].split('.')
            if spd[0] not in d.keys():
                d[spd[0]] = {}
            d[spd[0]][spd[1]] = val.strip(' \t')
        else:
            d[sp[0]] = val.strip(' \t')
    return d


def parseArray(content, field) -> []:
    tupls = re.findall(r'// @(' + field.lower() + '|' +
                       field.capitalize() + f')(.+)', content)
    return [tupl[1].strip(' \t') for tupl in tupls]


def parseMultiLine(content, field) -> str:
    tupls = re.compile(r'/\* @(' + field.lower() + '|' +
                       field.capitalize() + f')(.+?)\*/', re.DOTALL).findall(content)
    if len(tupls) == 0:
        return ''
    cn = re.sub(r'( \t|\t |\t)', '', tupls[0][1])
    return cn

###


def parseParam(content) -> []:
    tokens = re.findall(r'// @Param[ \t]+(.+)', content)
    p = []
    for token in tokens:
        token = token.split(' ')
        name, in_type, param_type, is_req, cmt = token[0], token[1], token[2], token[3], token[4:]
        cmt = ' '.join(cmt)
        p.append({
            "description": cmt.strip(' \t'),
            "in": in_type.strip(' \t'),
            "name": name.strip(' \t'),
            "required": True if is_req.strip(' \t').lower() == 'true' else False,
            "type": param_type.strip(' \t')
        })
    return p


def parseResp(content) -> []:
    # /* @ (~~~~~~) */
    tokens = re.compile(
        r'/\* @(success|Success|failed|Failed)(.+?)\*/', re.DOTALL).findall(content)
    c = []
    for token in tokens:
        token = token[1]
        sp = re.split(r'[ \t]', token)
        status, code = sp[0], sp[1]
        j = ' '.join(sp[2:]).strip(' \t')
        try:
            j = json.dumps(json.loads(j), indent=2)
        except:
            pass
        c.append({
            'status': status.strip(' \t'),
            'code': code.strip(' \t'),
            'json': j
        })
    return c

###


def parse_main(content):
    d = parsePlain(content)
    out = {}
    if 'basepath' in d:
        out['basePath'] = d.pop('basepath').strip(' \t')
    if 'host' in d:
        out['host'] = d.pop('host').strip(' \t')
    if 'swagger' in d:
        out['swagger'] = d.pop('swagger').strip(' \t')
    out['info'] = d
    if 'termsofservice' in d:
        d['termsOfService'] = d.pop('termsofservice')
    return out


def parse_ctrl(content, out_yaml):
    if 'paths' not in out_yaml:
        out_yaml['paths'] = {}

    contents = content.split('func ')
    for content in contents:
        try:
            # dels = ['param', 'success', 'failure', 'accept', 'produce']

            d = parsePlain(content)
            p = parseParam(content)
            c = parseResp(content)
            accept = parseArray(content, 'accept')
            desc = parseMultiLine(content, 'description')
            if desc == '' and 'description' in d:
                desc = d['description']

            # for de in dels:
            #     if de in d:
            #         del d[de]

            router = d['router']
            router, method = router.split(' ')
            method = method[1:-1].lower()
            oid = router.lower().replace('/', '-') + '-' + method
            oid = oid.replace('--', '-')

            if router not in out_yaml['paths']:
                out_yaml['paths'][router] = {}

            out_yaml['paths'][router][method] = {
                'summary': d['summary'],
                'description': desc,
                'consumes': accept,
                'produces': ['application/json'],
                'operationId': oid,
                'parameters': p,
                'responses': {code['code']: {'description': '```json\n%s\n```' % code['json']} for code in c},
            }
        except:
            continue


def main():
    main_file = sys.argv[1]
    all_files = []
    for root, _, files in os.walk('.'):
        for f in files:
            if f.split('.')[-1] == 'go':
                all_files.append(os.path.join(root, f))

    content = open(sys.argv[1], 'r', encoding='utf-8').read()
    out_yaml = parse_main(content)
    for f in all_files:
        content = open(f, 'r', encoding='utf-8').read()
        parse_ctrl(content, out_yaml)

    with open(sys.argv[2], 'w', encoding='utf-8') as f:
        yaml.dump(out_yaml, stream=f, encoding='utf-8', allow_unicode=True)


if __name__ == "__main__":
    main()

# python ./docs/parse_yaml.py main.go ./docs/api.yaml
