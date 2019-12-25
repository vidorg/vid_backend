import os
import sys
import re
import yaml
import json


def parsePlain(content) -> {}:
    """
    Parse plain @annotion string (2-lier)
    """
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
    """
    Parse plain @annotion array
    """
    tupls = re.findall(r'// @(' + field + '|' + field.lower() +
                       '|' + field.capitalize() + f')(.+)', content)
    return [tupl[1].strip(' \t') for tupl in tupls]


def parseMultiLine(content, field) -> str:
    """
    Parse plain @annotion multiline string
    """
    tupls = re.compile(r'/\* @(' + field.lower() + '|' +
                       field.capitalize() + f')(.+?)\*/', re.DOTALL).findall(content)
    if len(tupls) == 0:
        return ''
    cn = re.sub(r'( \t|\t |\t)', '', tupls[0][1])
    return cn


def parseParam(content) -> []:
    """
    Parse @Param name; in; type; required; description; enum; minLength; maxLength
    $name, $in, $type, $req, $desc, [$enum, $minLength, $maxLength]

    Ex: sex formData string false "new sex" enum(male, female, unknown)
    """
    tokens = re.findall(r'// @Param[ \t]+(.+)', content)
    p = []
    for token in tokens:
        token = re.split(r'[ \t]', token)
        name, in_type, param_type, is_req, others = token[0], token[1], token[2], token[3], token[4:]

        others = ' '.join(others)
        cmt = re.compile(r'.*"(.*)".*').findall(others)
        cmt = cmt[0].strip(' \t') if len(cmt) != 0 else ''
        enum = re.compile(r'.*enum\((.+?)\).*').findall(others)
        enum = [e.strip(' \t')
                for e in enum[0].split(',')] if len(enum) != 0 else []
        maxLength = re.compile(r'.*maxLength\((.+?)\).*').findall(others)
        maxLength = maxLength[0].strip(' \t') if len(maxLength) != 0 else ''
        minLength = re.compile(r'.*minLength\((.+?)\).*').findall(others)
        minLength = minLength[0].strip(' \t') if len(minLength) != 0 else ''

        pj = {
            "name": name.strip(' \t'),
            "description": cmt.strip(' \t'),
            "in": in_type.strip(' \t'),
            "required": True if is_req.strip(' \t').lower() == 'true' else False,
            "type": param_type.strip(' \t')
        }
        if len(enum) != 0:
            pj['enum'] = enum
        if minLength != '':
            pj['minLength'] = int(minLength)
            pj['description'] += '\n\n*Minimum length* : ' + minLength
        if maxLength != '':
            pj['maxLength'] = int(maxLength)
            if minLength == '':
                pj['description'] += '\n'
            pj['description'] += '\n*Maximin length* : ' + maxLength

        p.append(pj)
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


def parseErrorCode(content) -> str:
    """
    Parse UserDefined Error Code Message
    """
    ecs = parseArray(content, 'ErrorCode')
    codes, messages = [], []
    for ec in ecs:
        sp = re.split(r'[ \t]', ec)
        codes.append(int(sp[0].strip(' \t')))
        messages.append(' '.join(sp[1:]).strip(' \t'))
    md = '| Code | Message |\n| --- | --- |\n'
    for c, m in zip(codes, messages):
        md += '| %d | %s |\n' % (c, m)
    return md

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

            ecmd = parseErrorCode(content)
            desc += '\n\n' + ecmd + '\n\n'

            # for de in dels:
            #     if de in d:
            #         del d[de]

            router = d['router']
            sp = re.split(r'[ \t]', router)
            router, method = sp[0], sp[1]
            method = method[1:-1].lower()
            is_auth = len(sp) >= 3 and sp[2].lower() == '[auth]'

            oid = router.lower().replace('/', '-') + '-' + method
            oid = oid.replace('--', '-')

            yml = {
                'summary': d['summary'],
                'description': desc,
                'consumes': accept,
                'produces': ['application/json'],
                'operationId': oid,
                'parameters': p,
                'responses': {code['code']: {'description': '```json\n%s\n```' % code['json']} for code in c},
            }
            if is_auth:
                yml['security'] = [{
                    'basicAuth': '[]'
                }]

            if router not in out_yaml['paths']:
                out_yaml['paths'][router] = {}
            out_yaml['paths'][router][method] = yml
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
